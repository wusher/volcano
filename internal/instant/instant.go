// Package instant provides instant navigation with hover prefetching.
package instant

import "github.com/wusher/volcano/internal/minify"

// instantNavJSRaw is the unminified JavaScript code for instant navigation.
// It provides:
// - Hover-based link prefetching
// - Click interception for internal links
// - AJAX page loading with content replacement
// - History API integration
// - Theme state preservation
const instantNavJSRaw = `
(function() {
    'use strict';

    // Configuration
    const PREFETCH_DELAY = 65; // ms to wait before prefetching on hover
    const CONTENT_SELECTOR = '.content';
    const NAV_SELECTOR = '.tree-nav';
    const TOC_SELECTOR = '.toc-sidebar';
    const BREADCRUMBS_SELECTOR = '.breadcrumbs';
    const TITLE_SELECTOR = 'title';

    // Track prefetched URLs and pending prefetch
    const prefetched = new Set();
    let prefetchTimer = null;

    // Initialize instant navigation
    function init() {
        // Only run on pages that have the expected structure
        if (!document.querySelector(CONTENT_SELECTOR)) return;

        // Add event listeners
        document.addEventListener('mouseover', handleMouseOver, { passive: true });
        document.addEventListener('mouseout', handleMouseOut, { passive: true });
        document.addEventListener('click', handleClick);
        window.addEventListener('popstate', handlePopState);

        // Mark current page as "prefetched"
        prefetched.add(window.location.pathname);
    }

    // Handle mouseover for prefetching
    function handleMouseOver(e) {
        const link = e.target.closest('a');
        if (!link || !isInternalLink(link)) return;

        const href = link.getAttribute('href');
        if (!href || prefetched.has(href)) return;

        // Delay prefetch slightly to avoid false positives
        prefetchTimer = setTimeout(function() {
            prefetchPage(href);
        }, PREFETCH_DELAY);
    }

    // Handle mouseout to cancel pending prefetch
    function handleMouseOut(e) {
        if (prefetchTimer) {
            clearTimeout(prefetchTimer);
            prefetchTimer = null;
        }
    }

    // Check if link is internal (same origin, not a hash, not a download)
    function isInternalLink(link) {
        if (link.hostname !== window.location.hostname) return false;
        if (link.hasAttribute('download')) return false;
        if (link.hasAttribute('data-no-instant')) return false;
        const href = link.getAttribute('href');
        if (!href || href.startsWith('#')) return false;
        if (href.startsWith('mailto:') || href.startsWith('tel:') || href.startsWith('javascript:')) return false;
        // Skip same-page anchor links (e.g., /current-page/#section)
        if (link.pathname === window.location.pathname && link.hash) return false;
        return true;
    }

    // Prefetch a page using link prefetch
    function prefetchPage(href) {
        if (prefetched.has(href)) return;
        prefetched.add(href);

        const link = document.createElement('link');
        link.rel = 'prefetch';
        link.href = href;
        document.head.appendChild(link);
    }

    // Handle click on links
    function handleClick(e) {
        // Skip if modifier keys pressed
        if (e.ctrlKey || e.metaKey || e.shiftKey || e.altKey) return;
        if (e.button !== 0) return;

        const link = e.target.closest('a');
        if (!link || !isInternalLink(link)) return;

        const href = link.getAttribute('href');
        if (!href) return;

        e.preventDefault();
        navigateTo(href);
    }

    // Navigate to a new page via AJAX
    async function navigateTo(url, isPop) {
        try {
            // Show loading state
            document.body.style.cursor = 'progress';

            // Fetch the new page
            const response = await fetch(url);
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            const html = await response.text();

            // Parse the new page
            const parser = new DOMParser();
            const newDoc = parser.parseFromString(html, 'text/html');

            // Function to perform the actual DOM updates
            const performUpdate = () => {
                // Update URL FIRST so relative paths in new content resolve correctly
                // This must happen before updateContent to prevent broken images
                if (!isPop) {
                    history.pushState(null, '', url);
                }

                // Update page content
                updateContent(newDoc);

                // Scroll to top
                window.scrollTo({ top: 0, behavior: 'instant' });

                // Re-initialize page-specific features
                reinitialize();
            };

            // Use View Transitions API if available for smooth animations
            if (document.startViewTransition) {
                document.startViewTransition(performUpdate);
            } else {
                performUpdate();
            }

            // Reset cursor
            document.body.style.cursor = '';

        } catch (error) {
            console.error('Instant navigation failed:', error);
            // Fallback to normal navigation
            window.location.href = url;
        }
    }

    // Update page content from new document
    function updateContent(newDoc) {
        // Update title
        const newTitle = newDoc.querySelector(TITLE_SELECTOR);
        if (newTitle) {
            document.title = newTitle.textContent;
        }

        // Update main content
        const oldContent = document.querySelector(CONTENT_SELECTOR);
        const newContent = newDoc.querySelector(CONTENT_SELECTOR);
        if (oldContent && newContent) {
            oldContent.innerHTML = newContent.innerHTML;
        }

        // Update navigation (active states)
        const oldNav = document.querySelector(NAV_SELECTOR);
        const newNav = newDoc.querySelector(NAV_SELECTOR);
        if (oldNav && newNav) {
            oldNav.innerHTML = newNav.innerHTML;
        }

        // Update TOC
        const oldTOC = document.querySelector(TOC_SELECTOR);
        const newTOC = newDoc.querySelector(TOC_SELECTOR);
        if (oldTOC && newTOC) {
            oldTOC.innerHTML = newTOC.innerHTML;
            oldTOC.style.display = ''; // Restore display in case it was hidden
        } else if (oldTOC && !newTOC) {
            oldTOC.style.display = 'none';
        } else if (!oldTOC && newTOC) {
            // TOC appeared - insert it into the page
            const contentArea = document.querySelector('.content-area');
            if (contentArea) {
                contentArea.insertBefore(newTOC.cloneNode(true), contentArea.firstChild);
            }
        }

        // Update breadcrumbs
        const oldBreadcrumbs = document.querySelector(BREADCRUMBS_SELECTOR);
        const newBreadcrumbs = newDoc.querySelector(BREADCRUMBS_SELECTOR);
        if (oldBreadcrumbs && newBreadcrumbs) {
            oldBreadcrumbs.innerHTML = newBreadcrumbs.innerHTML;
        }
    }

    // Re-initialize features after content update
    function reinitialize() {
        // Dispatch custom event for user scripts
        document.dispatchEvent(new CustomEvent('instant:navigated', {
            detail: { url: window.location.href }
        }));
    }

    // Handle browser back/forward
    function handlePopState(e) {
        navigateTo(window.location.href, true);
    }

    // Start when DOM is ready
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }
})();
`

// InstantNavJS is the minified JavaScript code for instant navigation.
// It is initialized in init() to ensure proper package initialization order.
var InstantNavJS string

func init() {
	InstantNavJS = minify.JS(instantNavJSRaw)
}
