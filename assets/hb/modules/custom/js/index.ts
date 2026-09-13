// This script will be compiled into the JS bundle automatically.
import Collapse from "js/bootstrap/src/collapse";

(() => {
    function expandAccordionItem(id: string) {
        if (!id) return;
        const item = document.getElementById(id);
        if (!item || !item.classList.contains("accordion-item")) return;
        const collapse = item.querySelector(".accordion-collapse");
        if (!collapse) return;
        const bsCollapse = Collapse.getOrCreateInstance(collapse);
        if (!collapse.classList.contains("show")) {
            bsCollapse.show();
        }
    }

    function expandAccordionFromHash() {
        expandAccordionItem(window.location.hash.slice(1));
    }

    function initSamePageHashLinks() {
        document.addEventListener("click", (e) => {
            const target = e.target as Element | null;
            const link = target?.closest("a[href*='#']");
            if (!link) return;
            const url = new URL((link as HTMLAnchorElement).href, window.location.href);
            if (url.pathname !== window.location.pathname || !url.hash) return;
            expandAccordionItem(url.hash.slice(1));
        });
    }

    function initAnchorLinks() {
        document.querySelectorAll(".accordion-button .anchor-link").forEach((link) => {
            link.addEventListener("click", (e) => {
                e.stopPropagation();
            });
        });
    }

    if (document.readyState === "loading") {
        document.addEventListener("DOMContentLoaded", () => {
            expandAccordionFromHash();
            initAnchorLinks();
        });
    } else {
        expandAccordionFromHash();
        initAnchorLinks();
    }
    initSamePageHashLinks();
    window.addEventListener("hashchange", expandAccordionFromHash);
})()
