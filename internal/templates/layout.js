// Disable browser scroll restoration - we handle it manually
if ('scrollRestoration' in history) {
    history.scrollRestoration = 'manual';
}
// Scroll to top on page load (unless navigating to a hash)
// Use pageshow event to run after View Transitions complete
window.addEventListener('pageshow', function() {
    if (!window.location.hash) {
        window.scrollTo(0, 0);
    }
});

// Theme toggle
function toggleTheme() {
    const current = document.documentElement.getAttribute('data-theme');
    const next = current === 'dark' ? 'light' : 'dark';
    document.documentElement.setAttribute('data-theme', next);
    localStorage.setItem('theme', next);
    // Update browser theme-color meta tag
    const themeColor = next === 'dark' ? '#1a1a1a' : '#ffffff';
    const meta = document.querySelector('meta[name="theme-color"]');
    if (meta) meta.setAttribute('content', themeColor);
}

// Focus mode toggle (not persistent)
function toggleZenMode() {
    document.body.classList.toggle('zen-mode');
}

// Mobile drawer
function toggleDrawer() {
    document.body.classList.toggle('drawer-open');
}

function closeDrawer() {
    document.body.classList.remove('drawer-open');
}

// Mobile TOC toggle
function toggleMobileTOC() {
    document.body.classList.toggle('toc-open');
}

function closeMobileTOC() {
    document.body.classList.remove('toc-open');
}

// Close mobile TOC when clicking outside
document.addEventListener('click', function(e) {
    if (!document.body.classList.contains('toc-open')) return;
    const tocSidebar = document.querySelector('.toc-sidebar');
    const tocToggle = document.querySelector('.mobile-toc-toggle');
    if (tocSidebar && !tocSidebar.contains(e.target) && tocToggle && !tocToggle.contains(e.target)) {
        closeMobileTOC();
    }
});

// Close mobile TOC when clicking a TOC link
(function() {
    const tocSidebar = document.querySelector('.toc-sidebar');
    if (tocSidebar) {
        tocSidebar.addEventListener('click', function(e) {
            if (e.target.closest('a')) {
                closeMobileTOC();
            }
        });
    }
})();

// Tree navigation toggle (uses event delegation to survive instant nav)
document.querySelector('.tree-nav').addEventListener('click', function(e) {
    const toggle = e.target.closest('.folder-toggle');
    if (!toggle) return;
    e.preventDefault();
    const li = toggle.closest('li');
    li.classList.toggle('expanded');
});

// Close drawer on mobile when navigation link is clicked
document.querySelector('.tree-nav').addEventListener('click', function(e) {
    const link = e.target.closest('a.file-link, a.folder-link');
    if (!link) return;
    // Close the mobile drawer when a link is clicked
    closeDrawer();
});

// Expand path to current page
function expandActivePath() {
    document.querySelectorAll('.tree-nav a.active').forEach(function(active) {
        let parent = active.closest('li');
        while (parent) {
            parent.classList.add('expanded');
            parent = parent.parentElement.closest('li');
        }
    });
}
expandActivePath();
// Re-run after instant navigation
document.addEventListener('instant:navigated', expandActivePath);

// Scroll progress indicator
(function() {
    const progressBar = document.querySelector('.scroll-progress-bar');
    if (!progressBar) return;

    function updateProgress() {
        const scrollTop = window.scrollY;
        const docHeight = document.documentElement.scrollHeight - window.innerHeight;
        const progress = docHeight > 0 ? (scrollTop / docHeight) * 100 : 0;
        progressBar.style.width = progress + '%';
    }

    window.addEventListener('scroll', updateProgress, { passive: true });
    updateProgress();
})();

// Back to top button
(function() {
    const backToTop = document.querySelector('.back-to-top');
    if (!backToTop) return;

    const showThreshold = 300;

    function toggleBackToTop() {
        if (window.scrollY > showThreshold) {
            backToTop.hidden = false;
            backToTop.classList.add('visible');
        } else {
            backToTop.classList.remove('visible');
        }
    }

    backToTop.addEventListener('click', function() {
        window.scrollTo({ top: 0, behavior: 'smooth' });
    });

    window.addEventListener('scroll', toggleBackToTop, { passive: true });
})();

// Copy code button
function initializeCopyButtons() {
    document.querySelectorAll('.copy-button').forEach(function(button) {
        button.addEventListener('click', async function() {
            const code = this.parentElement.querySelector('code').textContent;
            try {
                await navigator.clipboard.writeText(code);
                this.classList.add('copied');
                this.setAttribute('aria-label', 'Copied!');
                setTimeout(function() {
                    button.classList.remove('copied');
                    button.setAttribute('aria-label', 'Copy code to clipboard');
                }, 2000);
            } catch (err) {
                console.error('Failed to copy:', err);
            }
        });
    });
}

// Initialize on page load
initializeCopyButtons();

// Reinitialize after instant navigation
document.addEventListener('instant:navigated', initializeCopyButtons);

// TOC scroll spy
(function() {
    let observer = null;

    function initializeTOCScrollSpy() {
        const toc = document.querySelector('.toc');
        if (!toc) return;

        // Disconnect previous observer if it exists
        if (observer) {
            observer.disconnect();
        }

        const headings = document.querySelectorAll('h2[id], h3[id], h4[id]');
        const tocLinks = toc.querySelectorAll('a');

        observer = new IntersectionObserver(function(entries) {
            entries.forEach(function(entry) {
                if (entry.isIntersecting) {
                    tocLinks.forEach(function(a) { a.classList.remove('active'); });
                    const link = toc.querySelector('a[href="#' + entry.target.id + '"]');
                    if (link) link.classList.add('active');
                }
            });
        }, { rootMargin: '-80px 0px -80% 0px' });

        headings.forEach(function(h) { observer.observe(h); });
    }

    // Initialize on page load
    initializeTOCScrollSpy();

    // Reinitialize after instant navigation
    document.addEventListener('instant:navigated', initializeTOCScrollSpy);
})();

// TOC smooth scroll with proper positioning
(function() {
    const toc = document.querySelector('.toc');
    if (!toc) return;

    const headerOffset = 80;

    // Use event delegation for TOC links
    toc.addEventListener('click', function(e) {
        const link = e.target.closest('a[href^="#"]');
        if (!link) return;

        e.preventDefault();
        const targetId = link.getAttribute('href').slice(1);
        const target = document.getElementById(targetId);
        if (!target) return;

        // Calculate position to put heading near top of viewport
        const targetPosition = target.getBoundingClientRect().top + window.scrollY - headerOffset;

        // Update active state immediately
        toc.querySelectorAll('a').forEach(function(a) { a.classList.remove('active'); });
        link.classList.add('active');

        // Use native smooth scroll (same as back-to-top button)
        // CSS handles prefers-reduced-motion via scroll-behavior: auto
        window.scrollTo({ top: targetPosition, behavior: 'smooth' });

        // Update URL hash
        history.pushState(null, null, '#' + targetId);
    });
})();

// Keyboard shortcuts
(function() {
    const baseURL = window.VOLCANO_BASE_URL || '';
    const homeURL = baseURL ? baseURL + '/' : '/';
    const shortcuts = {
        't': function() { toggleTheme(); },
        'z': function() { toggleZenMode(); },
        'h': function() { window.location.href = homeURL; },
        '?': function() { showShortcutsModal(); }
    };


    // Update page navigation shortcuts (p/n keys)
    function updatePageNavShortcuts() {
        const prevLink = document.querySelector('.page-nav-prev');
        const nextLink = document.querySelector('.page-nav-next');

        // Remove old shortcuts
        delete shortcuts['p'];
        delete shortcuts['n'];

        // Add new shortcuts if links exist
        if (prevLink) shortcuts['p'] = function() { window.location.href = prevLink.href; };
        if (nextLink) shortcuts['n'] = function() { window.location.href = nextLink.href; };
    }

    // Initialize on page load
    updatePageNavShortcuts();

    // Update after instant navigation
    document.addEventListener('instant:navigated', updatePageNavShortcuts);

    // Global keydown handler
    document.addEventListener('keydown', function(e) {
        // Skip if in input/textarea
        if (['INPUT', 'TEXTAREA'].includes(e.target.tagName)) {
            if (e.key === 'Escape') {
                e.target.blur();
                closeDrawer();
            }
            return;
        }

        // Close modal/zen mode on Escape
        if (e.key === 'Escape') {
            closeShortcutsModal();
            closeDrawer();
            closeMobileTOC();
            document.body.classList.remove('zen-mode');
            return;
        }

        const handler = shortcuts[e.key];
        if (handler) {
            e.preventDefault();
            handler();
        }
    });
})();

// Shortcuts modal
function showShortcutsModal() {
    const modal = document.getElementById('shortcuts-modal');
    if (modal && modal.showModal) modal.showModal();
}

function closeShortcutsModal() {
    const modal = document.getElementById('shortcuts-modal');
    if (modal && modal.close) modal.close();
}
