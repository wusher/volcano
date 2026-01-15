package search

// GenerateSearchJS returns the JavaScript code for the command palette.
// This is loaded on-demand when user first presses Cmd+K.
func GenerateSearchJS(baseURL string) string {
	return `(function() {
    const baseURL = '` + baseURL + `';
    let searchIndex = null;
    let searchLoading = false;
    let selectedIndex = -1;

    // Inject HTML
    const html = '<div id="command-palette" class="command-palette">' +
        '<div class="command-palette-backdrop"></div>' +
        '<div class="command-palette-modal" role="dialog" aria-label="Search">' +
            '<div class="command-palette-header">' +
                '<input type="text" id="command-palette-input" placeholder="Search..." autocomplete="off" spellcheck="false">' +
                '<kbd class="command-palette-hint">esc</kbd>' +
            '</div>' +
            '<div class="command-palette-results" id="command-palette-results">' +
                '<div class="command-palette-empty">Type to search...</div>' +
            '</div>' +
        '</div>' +
    '</div>';
    document.body.insertAdjacentHTML('beforeend', html);

    async function loadSearchIndex() {
        if (searchIndex || searchLoading) return searchIndex;
        searchLoading = true;
        try {
            const url = baseURL ? baseURL + '/search-index.json' : '/search-index.json';
            const res = await fetch(url);
            searchIndex = await res.json();
        } catch (e) {
            console.error('Failed to load search index:', e);
        }
        searchLoading = false;
        return searchIndex;
    }

    const palette = document.getElementById('command-palette');
    const input = document.getElementById('command-palette-input');
    const results = document.getElementById('command-palette-results');

    function openCommandPalette() {
        palette.classList.add('open');
        document.body.classList.add('command-palette-open');
        input.value = '';
        selectedIndex = -1;
        results.innerHTML = '<div class="command-palette-empty">Type to search...</div>';
        input.focus();
        loadSearchIndex();
    }

    function closeCommandPalette() {
        palette.classList.remove('open');
        document.body.classList.remove('command-palette-open');
    }

    window.openCommandPalette = openCommandPalette;
    window.closeCommandPalette = closeCommandPalette;

    document.addEventListener('keydown', function(e) {
        if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
            e.preventDefault();
            openCommandPalette();
            return;
        }
        if (e.key === 'Escape' && palette.classList.contains('open')) {
            e.preventDefault();
            closeCommandPalette();
            return;
        }
    });

    document.querySelector('.command-palette-backdrop').addEventListener('click', closeCommandPalette);

    function escapeHtml(s) {
        return s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
    }

    function doSearch() {
        const query = input.value.trim().toLowerCase();
        if (!searchIndex || !query) {
            results.innerHTML = '<div class="command-palette-empty">Type to search...</div>';
            selectedIndex = -1;
            return;
        }

        const terms = query.split(/\s+/).filter(t => t.length > 0);
        const matches = [];
        for (const page of searchIndex.pages) {
            const pageTitle = page.title.toLowerCase();
            const pagePath = page.url.toLowerCase().replace(/\//g, ' ');
            const pageText = pageTitle + ' ' + pagePath;
            // Page matches if all terms are in title or path
            if (terms.every(t => pageText.includes(t))) {
                matches.push({ type: 'page', title: page.title, url: page.url, snippet: '' });
            }
            // Heading matches if all terms are in heading, title, or path combined
            for (const h of (page.headings || [])) {
                const combined = h.text.toLowerCase() + ' ' + pageText;
                if (terms.every(t => combined.includes(t))) {
                    matches.push({ type: 'h' + h.level, title: h.text, url: page.url + '#' + h.anchor, snippet: page.title });
                }
            }
        }

        if (matches.length === 0) {
            results.innerHTML = '<div class="command-palette-empty">No results found</div>';
            selectedIndex = -1;
            return;
        }

        const limited = matches.slice(0, 10);
        const prefixedBaseURL = baseURL || '';
        results.innerHTML = limited.map(function(m, i) {
            return '<a href="' + prefixedBaseURL + m.url + '" class="command-palette-result' + (i === 0 ? ' selected' : '') + '" data-index="' + i + '">' +
                '<span class="result-type">' + m.type + '</span>' +
                '<span class="result-title">' + escapeHtml(m.title) + '</span>' +
                (m.snippet ? '<span class="result-snippet">' + escapeHtml(m.snippet) + '</span>' : '') +
            '</a>';
        }).join('');
        selectedIndex = 0;
    }

    input.addEventListener('input', doSearch);

    input.addEventListener('keydown', function(e) {
        const items = results.querySelectorAll('.command-palette-result');
        if (items.length === 0) return;

        if (e.key === 'ArrowDown') {
            e.preventDefault();
            selectedIndex = Math.min(selectedIndex + 1, items.length - 1);
            updateSelection(items);
        } else if (e.key === 'ArrowUp') {
            e.preventDefault();
            selectedIndex = Math.max(selectedIndex - 1, 0);
            updateSelection(items);
        } else if (e.key === 'Enter' && selectedIndex >= 0) {
            e.preventDefault();
            items[selectedIndex].click();
        }
    });

    function updateSelection(items) {
        items.forEach(function(item, i) {
            item.classList.toggle('selected', i === selectedIndex);
        });
    }

    results.addEventListener('click', function(e) {
        if (e.target.closest('.command-palette-result')) {
            closeCommandPalette();
        }
    });

    // Listen for open-search event (from mobile button)
    window.addEventListener('open-search', openCommandPalette);

    // Open immediately since user pressed Cmd+K to load this
    openCommandPalette();
})();`
}
