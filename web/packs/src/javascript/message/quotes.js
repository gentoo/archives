$(function() {

    var message = document.querySelector("pre.ag-message-content").innerHTML;
    var lines = message.split("\n");
    var convertedLines = [];
    var inQuote = false;

    lines.forEach(function (line) {
        if(line.trimLeft().startsWith("&gt;") || line.trimLeft().startsWith(">")){
            if(inQuote){
                convertedLines.push(line);
            }else{
                convertedLines.push('<div class="ag-quote">');
                convertedLines.push(line);
                inQuote = true;
            }
        }else{
            if(inQuote){
                convertedLines.push('</div>');
                convertedLines.push(line);
                inQuote = false;
            }else{
                convertedLines.push(line);
            }
        }
    })

    document.querySelector("pre.ag-message-content").innerHTML = convertedLines.join('\n');

    console.log('Done');
})