import hljs from 'highlight.js/lib/core';
import diff from 'highlight.js/lib/languages/diff';

$(function() {

    console.log("Format git!");

    let message = document.querySelector("pre.ag-message-content").innerHTML;
    let gitDiffRegex = new RegExp('diff --git a\\/(.*)\\sb\\/(.*)\\nindex\\s([a-zA-Z0-9]*)\\.\\.([a-zA-Z0-9]*)\\s([a-zA-Z0-9]*)\\n---\\sa\\/(.*)\\n\\+\\+\\+\\sb\\/(.*)\\n');
    let isGitDiff = (message.search(gitDiffRegex) != -1);

    if(isGitDiff){
        hljs.registerLanguage('patch', diff);
        hljs.registerLanguage('diff', diff);

        document.querySelector("pre.ag-message-content").innerHTML = '<code class="language-patch">' + message + '</code>';
        document.querySelectorAll('pre.ag-message-content code').forEach((block) => {
            hljs.highlightBlock(block);
        });
    }
})


