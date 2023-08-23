pragma solidity >=0.4.21 <0.7.0;
pragma experimental ABIEncoderV2;

//(Role_Manange_Contract)

contract RMC{

	struct RolePolicy{
		string Role;
		string[] Attributes;  //name_lhz
	}
	
	mapping (string => RolePolicy) public PolicyList;
	address public owner;
	constructor() public {
        owner = msg.sender;
	}
	
	function writePolicy(string memory p_1, string memory p_2, string memory p_3, string memory role, string memory id) public {
		string[] memory attris = new string[](5);
		RolePolicy memory rPolicy = RolePolicy("",attris);
		rPolicy.Attributes = attris;
		rPolicy.Attributes[0]=p_1;
		rPolicy.Attributes[1]=p_2;
		rPolicy.Attributes[2]=p_3;
		rPolicy.Role = role; 
		PolicyList[id] = rPolicy;
	}
	
	function writePolicyT(string[] memory attributes, string memory role, string memory id ) public {
		//uint len = attributes.length;
	    RolePolicy memory rPolicy = RolePolicy("",attributes);
		rPolicy.Attributes = attributes;
		rPolicy.Role = role; 
		PolicyList[id] = rPolicy;
	}
	
	
	function assignRole(string memory policy_id, string memory p_1, string memory p_2, string memory p_3) public view returns (string memory){
		RolePolicy memory rPolicy = PolicyList[policy_id];
		
		if (keccak256(abi.encodePacked(rPolicy.Attributes[0])) == keccak256(abi.encodePacked(p_1))){
			if(keccak256(abi.encodePacked(rPolicy.Attributes[1])) == keccak256(abi.encodePacked(p_2))){
				if(keccak256(abi.encodePacked(rPolicy.Attributes[2])) == keccak256(abi.encodePacked(p_3))){
					return rPolicy.Role;
				}
			}
		}
		return "NORole";
	}
	
	
	//数组类型的参数
	function assignRoleT(string[] memory attributes, string memory policy_id) public view returns (string memory){
	
		RolePolicy memory rPolicy = PolicyList[policy_id];
		uint count = 0;
		
		for (uint i = 0; i <= attributes.length; i++) {		
			for (uint j = 0; j<rPolicy.Attributes.length; j++){
				if (keccak256(abi.encodePacked(rPolicy.Attributes[j])) == keccak256(abi.encodePacked(attributes[i]))){
					count = count +1;	
				}
			}
		}
		if (count == attributes.length){
			return rPolicy.Role;
		}
		return "NORole";
	}
	
	
	

}

//客户端调用规则：
//先完成数据的set操作之后
//随后调用getACL查看
//最后调用setFileRoleAccess设定文件对应权限
