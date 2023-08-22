pragma solidity >=0.4.21 <0.7.0;
//pragma experimental ABIEncoderV2;

//(Role_Permission_Manange_Contract)

contract RPMC{

    struct Permission {
		 uint Individual;  //个人信息
		 uint Physiology;  //生理数据
		 uint Diagnosis;   //诊断信息
	}
	
	enum Role {Doctors, Pharmacists, Academics}
	
	Permission[] public Role_access;
	
	address public owner;
	
	mapping (uint => string) type_map;
	
	mapping(string => string) Files_Role_Access;
	
	
	constructor() public {
        owner = msg.sender;
		
		Permission memory permission =  Permission(0,0,0);
		Role_access.push(permission); 
		Role_access.push(permission); 
		Role_access.push(permission); 
		
		type_map[0] ="0";
		type_map[1] ="1";
    }
	
	
    function SetRoleAcl(uint _Individual,uint _Physiology, uint _Diagnosis, string memory role) public {
	    
		Permission memory permission = Permission(_Individual,_Physiology,_Diagnosis);
		
		if (keccak256(abi.encodePacked(role)) == keccak256(abi.encodePacked("Doctors"))){ 
			Role_access[uint(Role.Doctors)] = permission;	
		}
		
		if (keccak256(abi.encodePacked(role)) == keccak256(abi.encodePacked("Pharmacists"))){
		
			Role_access[uint(Role.Pharmacists)] = permission;
		}
		
		if (keccak256(abi.encodePacked(role)) == keccak256(abi.encodePacked("Academics"))){
		
			Role_access[uint(Role.Academics)] = permission;
			
		}		
	
	}
	
	
	function getRoleAcl(string memory role) public view returns (uint I, uint P, uint D) {
	
		if (keccak256(abi.encodePacked(role)) == keccak256(abi.encodePacked("Doctors"))){ 
			
			return (Role_access[uint(Role.Doctors)].Individual,Role_access[uint(Role.Doctors)].Physiology, Role_access[uint(Role.Doctors)].Diagnosis);  
		}
		
		if (keccak256(abi.encodePacked(role)) == keccak256(abi.encodePacked("Pharmacists"))){
			
			return (Role_access[uint(Role.Pharmacists)].Individual,Role_access[uint(Role.Pharmacists)].Physiology, Role_access[uint(Role.Pharmacists)].Diagnosis);
		}
	
		if (keccak256(abi.encodePacked(role)) == keccak256(abi.encodePacked("Academics"))){
			
			return (Role_access[uint(Role.Academics)].Individual,Role_access[uint(Role.Academics)].Physiology, Role_access[uint(Role.Academics)].Diagnosis);
		}
	}
	
	
	
	function setFileRoleAccess(string memory fileID) public {
		
		uint key = Role_access[uint(Role.Doctors)].Individual;  
		
		string memory str_1 = type_map[key]; 
		
		key = Role_access[uint(Role.Doctors)].Physiology;
		
		string memory str_2 = type_map[key];
		
		key = Role_access[uint(Role.Doctors)].Diagnosis;
		
		string memory str_3 = type_map[key];
		
		string memory tmp_str = strConcat("",str_1);
		
		tmp_str = strConcat(tmp_str,str_2);
		
		tmp_str = strConcat(tmp_str,str_3);
		
		
		
		key = Role_access[uint(Role.Pharmacists)].Individual;  
		
		str_1 = type_map[key];             
		
		key = Role_access[uint(Role.Pharmacists)].Physiology;
		
		str_2 =type_map[key]; 
		
		key = Role_access[uint(Role.Pharmacists)].Diagnosis;
		
		str_3 = type_map[key];
		
		tmp_str = strConcat(tmp_str,":");
		
		tmp_str = strConcat(tmp_str,str_1);
		
		tmp_str = strConcat(tmp_str,str_2);
		
		tmp_str = strConcat(tmp_str,str_3);
		
		
	
		key = Role_access[uint(Role.Academics)].Individual;
	
		str_1 = type_map[key]; 
		
		key = Role_access[uint(Role.Academics)].Physiology;
		
		str_2 =type_map[key]; 
		
		key = Role_access[uint(Role.Academics)].Diagnosis;
		
		str_3 = type_map[key];
		
		tmp_str = strConcat(tmp_str,":");
		
		tmp_str = strConcat(tmp_str,str_1);
		
		tmp_str = strConcat(tmp_str,str_2);
		
		tmp_str = strConcat(tmp_str,str_3);
		
		
		Files_Role_Access[fileID] = tmp_str;
		
	}
	
	function getFileRoleAccess(string memory fileID) public view returns (string memory ret) {
		
		return (Files_Role_Access[fileID]);
	}
	
	
	
	function strConcat(string memory _a, string memory _b) public view returns (string memory){
        bytes memory _ba = bytes(_a);
        bytes memory _bb = bytes(_b);
        string memory ret = new string(_ba.length + _bb.length);
        bytes memory bret = bytes(ret);
        uint k = 0;
        uint i = 0;
        for (i = 0; i < _ba.length; i++)bret[k++] = _ba[i];
        for (i = 0; i < _bb.length; i++) bret[k++] = _bb[i];
        return string(ret);
    } 
	
}

//客户端调用规则：
//先完成数据的set操作之后
//随后调用getACL查看
//最后调用setFileRoleAccess设定文件对应权限
