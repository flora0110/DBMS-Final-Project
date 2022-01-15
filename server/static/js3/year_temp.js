var page=0;
var len;
var whichYear=-1;
var start_2021=0,start_2020=0;
function next() {
    page+=50;
    //alert(page);
    var button1 = document.getElementById('pre');
    if(page<=0) {
		    button1.disabled = true;
    }
    else{
        button1.disabled = false;
    }
    var button2 = document.getElementById('next');
    if(page>=(len-50) || (whichYear==2021 && page>(start_2020-50)) ) {
		    button2.disabled = true;
    }
    else{
        button2.disabled = false;
    }
    console.log(page);
    drawTable();
}
function pre() {
    page-=50;
    //alert(page);
    var button1 = document.getElementById('pre');
    if(page<=0 || (whichYear==2020 && page<(start_2020)) ) {
		    button1.disabled = true;
    }
    else{
        button1.disabled = false;
    }
    var button2 = document.getElementById('next');
    if(page>=(len-50) || (whichYear==2021 && page>(start_2020-50)) ) {
		    button2.disabled = true;
    }
    else{
        button2.disabled = false;
    }
    console.log(page);
    drawTable();
}
var newdata =[];

//console.log("aa");
init({{ data|safe }});
drawTable();

function init(data){
    //console.log("bb");
    len = data.length;
    var k=0;
    for (var i =0 ; i < len; i++) {
        if(data[i][1]==2021) {
            newdata.push(data[i]);
        }
    }
    start_2020=newdata.length+1;
    //console.log(start_2020);
    for (var i =0 ; i < len; i++) {
        if(data[i][1]==2021) {break;}
        newdata.push(data[i]);

    }
    //console.log(newdata);
}
function drawTable(){
    var forTable = document.querySelector('tbody');
    for (var i =page ; i < page+50; i++) {
        if(i>len) break;
        if(whichYear==2021 && i>start_2020) break;
        if(newdata[i][1]==whichYear || whichYear==-1){
            for(var j=0;j<5;j++){
                if(newdata[i][j]=="NULL") {newdata[i][j]="";}
            }

            forTable.innerHTML +=
                '<tr>' +
                '<th scope="row" class="h6 text-black-50" >'+(i+1)+'</th>'+
                '<td>' + newdata[i][1] + '</td>' +
                '<td>' + newdata[i][2] + '</td>' +
                '<td>' + newdata[i][0] + '</td>' +
                '<td>' + newdata[i][3] + '</td>' +
                '<td>' + newdata[i][4] + '</td>' +
            '</tr>';
        }


    }
}

$(document).ready(function(){
  $("#next").click(function(){
    $("#test").empty();
    next();
  });
});
$(document).ready(function(){
  $("#pre").click(function(){
    $("#test").empty();
    pre();
  });
});
function changeYear(data){
    if(data==2020) {
        $(document).ready(function(){
            $("#test").empty();
            page=start_2020;
            whichYear=2020;
            var button1 = document.getElementById('pre');
            button1.disabled = true;
            drawTable();
        });
    }
    if(data==2021) {
        $(document).ready(function(){
            $("#test").empty();
            page=start_2021;
            whichYear=2021;
            var button1 = document.getElementById('pre');
            button1.disabled = true;

            drawTable();
        });
    }
    if(data=="all") {
        $(document).ready(function(){
            $("#test").empty();
            page=0;
            var button1 = document.getElementById('pre');
            button1.disabled = true;
            whichYear=-1;
            drawTable();
        });
    }
    console.log(data);
}