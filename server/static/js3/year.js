console.log("eee")
$(document).ready(function() {
    var forTable = document.querySelector('.for-table tbody');
    for (var i =0 ; i < 5; i++) {
            forTable.innerHTML +=
                '<tr>' +
                    '<td>' + {{ data[i][0] }} + '</td>' +
                    '<td>' + {{ data[i][1] }} + '</td>' +
                '</tr>';
        }
});
