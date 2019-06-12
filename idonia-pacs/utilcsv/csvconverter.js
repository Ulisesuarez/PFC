let fs = require('fs');

let args = process.argv;
let type = args[2] || 'text';
let arr = []; 
let bufferString; 


  fs.readFile('./tags2.csv',function (err,data) {

  if (err) {
    return console.log(err);
  }

  //Convert and store csv information into a buffer. 
  bufferString = data.toString(); 

  //Store information for each individual person in an array index. Split it by every newline in the csv file. 
  arr = bufferString.split('\n'); 
  let jsonObj = {};
  let headers = arr[0].split(';');
  for(let i = 1; i < arr.length; i++) {
    let data = arr[i].split(';');
    let obj = {};
    for(let j = 1; j < data.length; j++) {
        console.log()
       if (headers[j].trim() ==='VR'){
        obj['VR'] = data[j].trim().split('or');
       } else {
        obj[headers[j].trim()] = data[j].trim();
       }
    }
    jsonObj[data[0].trim()] = obj;
  }  
  //console.log(JSON.stringify(jsonObj));
  fs.writeFile("./tags.json", JSON.stringify(jsonObj), function(err) {
    if(err) {
        return console.log(err);
    }

    console.log("The file was saved!");
}); 
});
