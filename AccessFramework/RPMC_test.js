var Web3 = require('web3');
if (typeof web3 !== 'undefined') {
	web3 = new Web3(web3.currentProvider);
} else {
	// set the provider you want from Web3.providers
	web3 = new Web3(new Web3.providers.HttpProvider("http://127.0.0.1:8545"));
}

var fs = require('fs');

var RPMC_JSON = JSON.parse(fs.readFileSync("./build/contracts/RPMC.json"));

//console.log(RPMC_JSON);
var Mycontract = new web3.eth.Contract(RPMC_JSON.abi);

if (typeof Mycontract == 'undefined'){
	console.log("Mycontract undefined");
}
Mycontract.options.address = "0x52084F1D0Bcf872447Adc982CDd67F0a7DAbC4C9";  //合约地址要注意是否变化 然后实时更改
/////////////////////////////////////////////////////////////////////////////////////////服务器配置及初始化参数


function setFileAclByTransactions(fileId,doc_acl,pha_acl,aca_acl){   //各个角色权限的设置与规划
  
  var d_arr = doc_acl.split('_');
  var p_arr = pha_acl.split('_');
  var a_arr = aca_acl.split('_');
  console.log("Permission Control Matrix:");
  console.log(d_arr);
  console.log(p_arr);
  console.log(a_arr);
  
	Mycontract.methods.SetRoleAcl(d_arr[0], d_arr[1], d_arr[2], "Doctors").send({ from:"0x3F93D8eB7A099BB8e0404f723517779db5420cF9",gas:123456}).on('receipt', (data) => {
	 console.log(data);
  });
	Mycontract.methods.SetRoleAcl(p_arr[0], p_arr[1], p_arr[2], "Pharmacists").send({ from:"0x3F93D8eB7A099BB8e0404f723517779db5420cF9",gas:123456}).on('receipt', (data) => {
	console.log(data);
  });
	Mycontract.methods.SetRoleAcl(a_arr[0], a_arr[1], a_arr[2], "Academics").send({ from:"0x3F93D8eB7A099BB8e0404f723517779db5420cF9",gas:123456}).on('receipt', (data) => {
	console.log(data);
  });
	Mycontract.methods.setFileRoleAccess(fileId).send({ from:"0x3F93D8eB7A099BB8e0404f723517779db5420cF9",gas:123456}).on('receipt', (data) => {
	console.log(data);
	ret = data;
  });
}

var acl_role;
function getFileRoleAcl(fileId){  //fileId: "01234546789"
   Mycontract.methods.getFileRoleAccess(fileId).call().then(function(result){
    console.log("function---getFileRoleAccess");
	console.log("AcList");
	console.log(result);
	acl_role = result;
   });   
}



console.log("Server Initing.....");
const http = require('http');
const ip = '127.0.0.1';
const port = 3000;

http.createServer((req, res) => {
	console.log("Server Working....");
	var url = require('url');
	var params = url.parse(req.url, true).query;
    console.log("operartion：" + params.op);   //解析参数res.write("operartion：" + params.op);
    console.log("fileID：" + params.fileId); 
	var fileId = params.fileId;
	if(params.op == "get"){
		getFileRoleAcl(fileId);
		if (typeof acl_role == 'undefined'||acl_role ==""){
			res.end("error");
			console.log(acl_role);
		}else{
			res.end(acl_role);
			console.log("sending:"+acl_role);
		}
	}
	if(params.op == "set"){
		var doc_acl = params.d;
		var pha_acl = params.p;
		var aca_acl = params.a;
		setFileAclByTransactions(fileId,doc_acl,pha_acl,aca_acl);
		res.end("transaction sending...");
	}
	if (params.op == "current"){
		Mycontract.methods.getRoleAcl("Doctors").call().then(function(result){
			console.log("Doctors");
			console.log(result.I);
			console.log(result.P);
			console.log(result.D);
		}); //调用读数据合约
		Mycontract.methods.getRoleAcl("Pharmacists").call().then(function(result){
			console.log("Pharmacists");
			console.log(result.I);
			console.log(result.P);
			console.log(result.D);
		 });   
 
		 Mycontract.methods.getFileRoleAccess("01234546789").call().then(function(result){
			console.log("AcList");
			console.log(result);
		 }); 
	}
    //res.write("fileID：" + params.fileID); res.end(acl_role);
  }).listen(port, ip);
  
console.log(`server has started at ${ip}:${port}`);
 
 
 
 
 
 

//调用写数据的合约
//console.log(Mycontract);
// let instance = await RPMC.deployed("hello")
// instance.SetRoleAcl(1,1,1,"Doctors") 
// instance.SetRoleAcl(0,1,0,"Pharmacists")
// instance.SetRoleAcl(0,1,0,"Academics")
// instance.setFileRoleAccess("01234546789")
// instance.getFileRoleAccess("01234546789")
//解决Gas Out 错误需要指出 gas 参数