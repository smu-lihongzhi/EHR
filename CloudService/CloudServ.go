// CloudServ project main.go
package main

import (
	"os"
	"fmt"
	"io"
	"net/http"
	"math/rand"
	"time"
	"log"
	"strconv"
	"encoding/json"
	"strings"
)

type individual struct {
	Name string `json:"name"`
	Age string  `json:"age"`
	IDNum string `json:"IDNum"`
	Address string `json:"address"`
	Job string `json:"job"`
	Extra string `json:"extra"`
}

type physiology struct {
	HeartRate string `json:"heartRate"`
	Pulse string `json:"pulse"`
	BloodPres string `json:"bloodPres"`
	Extra string `json:"extra"`
}

type diagnosis struct{ //对应JSON字段的首字母要大写
	Content string `json:"content"`  
	Author string `json:"author"`
	Time int64 `json:"time"`
}

type UserData struct {
	Individual individual `json:"individual"`
	Physiology physiology `json:"physiology"`
	Diagnosis  diagnosis  `json:"diagnosis"`
}

var token string  //与设备端通信的token
var token_map map[string]string //token有效性

////////////////////////////////////////////////////////////////////////////
type Right struct{
	Role string    `json:"role"`
	Individual int `json:"individual"`
	Physiology int `json:"physiology"`
	Diagnosis  int `json:"diagnosis"`
}

type FilePermission struct{
	FileId string        `json:"fileId"`
	Rights [3]Right       `json:"rights"`
}
///////////////////////////////////////////////////////////////数据权限的定义


func initUserData(heart_rate string, blood_press string, pulse string) (UserData){
	
   var usr_data UserData
   usr_data.Individual.Address = "An Hui, Chizhou"   //这部分信息可由医院的就诊卡提供
   usr_data.Individual.Age = "35"
   usr_data.Individual.IDNum ="23455555"  //如就诊卡的卡号
   usr_data.Individual.Job = "teacher"
   usr_data.Individual.Name = "Hongzhi Li"
   usr_data.Individual.Extra = "xxxxxxxxx"
   usr_data.Physiology.HeartRate = heart_rate   //心率
   usr_data.Physiology.BloodPres = blood_press  //血压
   usr_data.Physiology.Pulse = pulse   //脉搏数据
   usr_data.Physiology.Extra ="xxxxxxxxx"
   //模拟的来自医生的医疗诊断数据
   usr_data.Diagnosis.Author = "black bill"
   usr_data.Diagnosis.Content = "Bronchoscopy rigid Bronchoscopy has been replaced by fiberoptic bronchoscopy" 
   usr_data.Diagnosis.Time = time.Now().Unix()
	
   return usr_data
}


func randomString(size int, kind int) string {  //随机的方式生成硬件设备的Token
    ikind, kinds, rsbytes := kind, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)  //随机数生成
    isAll := kind > 2 || kind < 0
    rand.Seed(time.Now().UnixNano())
    for i := 0; i < size; i++ {
        if isAll { // random ikind
            ikind = rand.Intn(3)
        }
        scope, base := kinds[ikind][0], kinds[ikind][1]
        rsbytes[i] = uint8(base + rand.Intn(scope))
    }
    return string(rsbytes)
}


func obtainToken(w http.ResponseWriter, r *http.Request){   //设定通信
	//fmt.Println(randomString(6, 0))
	params := r.URL.Query()
	key := params["key"][0]  //可以由请求参数获取设备节点的标识
	_,ok := token_map[key]
	if ok{
		tk := token_map[key]
		fmt.Fprintf(w, tk)
	}else{
		token = randomString(6, 0)
		token_map[key] = token 
		fmt.Fprintf(w, token)
	}
	fmt.Println("token is:",token) //输出false
}



func sendEMR(w http.ResponseWriter, r *http.Request){  //由设备终端调用,写医疗诊断数据EMC,url参数包括token、heart_rate、blood_press等
	params := r.URL.Query()
	tk := params["token"][0]   //请求URL中自带的通信Token 
	ok := false
	
	for _, v := range token_map {
        //fmt.Printf("Key:%v Value: %s\n", k, v)
        if v == tk {
        	  ok = true
        }
    }
	
	if ok {
		heart_rate := params["heart_rate"][0]
		blood_press := params["blood_press"][0]
		pulse := params["pulse"][0]
		usr_data := initUserData(heart_rate,blood_press,pulse)
		jsonBytes, err := json.Marshal(usr_data)
		if err != nil{
			fmt.Println(err)
			return
		}
		jsonString := string(jsonBytes)
    		fmt.Println(jsonString)
		
		file_name := "D:\\Truffle_test\\user_files\\user_data_"+tk+"_"+usr_data.Individual.IDNum+".js" //最好在文件名加入就诊卡号以标识
		var file *os.File
		var flag bool
		_, err = os.Stat(file_name)

	    if err == nil{
	        flag = true
	    }
		if os.IsNotExist(err){  //文件不存在
			flag = false
		}
		
		if flag == false {
			file, err = os.Create(file_name)
			if err != nil{
				fmt.Println("创建失败",err)
				return
			}
			
		}else{
			file,err = os.OpenFile(file_name,os.O_RDWR,0666)
			if err!=nil{
				fmt.Println("打开文件错误：",err)
				return
			}
		}
		defer file.Close()
		n,err:=io.WriteString(file,jsonString)
		if err != nil {
			fmt.Println("写入错误：",err)
			return
		}
		fmt.Println("写入成功：n=",n)
		//fmt.Fprintf(w, jsonString) 
		fileId := tk+"_"+usr_data.Individual.IDNum
		fmt.Fprintf(w,fileId)
		
	}else{
		fmt.Fprintf(w, "Token Error!")
		return
	}
	
}

func sendEMRRights(w http.ResponseWriter, r *http.Request){  //写医疗诊断数据EMC的访问控制文件,url参数包括d_r、p_r、a_r等
	
	var fileRights FilePermission
    var dRights Right
    var pRights Right
    var aRights Right
	params := r.URL.Query()

	d_rights := strings.Split(params["d_r"][0],"_")  //doctor rights 如:1_1_1
    p_rights := strings.Split(params["p_r"][0],"_")  //pharmacist rights 如:0_1_1 
	a_rights := strings.Split(params["a_r"][0],"_")  //academics rights 如:0_0_1
	
    dRights.Role = "Doctors"
    dRights.Individual,_ = strconv.Atoi(d_rights[0])
    dRights.Physiology,_ = strconv.Atoi(d_rights[1])
    dRights.Diagnosis,_ = strconv.Atoi(d_rights[2])
    
    pRights.Role = "Pharmacists"
    pRights.Individual,_ = strconv.Atoi(p_rights[0])
    pRights.Physiology,_ = strconv.Atoi(p_rights[1])
    pRights.Diagnosis,_ = strconv.Atoi(p_rights[2])
    
    aRights.Role = "Academics"
    aRights.Individual,_ = strconv.Atoi(a_rights[0])
    aRights.Physiology,_ = strconv.Atoi(a_rights[1])
    aRights.Diagnosis,_ = strconv.Atoi(a_rights[2])
    
    fileRights.FileId = params["fileId"][0]
	fileRights.Rights[0]= dRights
	fileRights.Rights[1]= pRights
	fileRights.Rights[2]= aRights
	
	jsonBytes, err := json.Marshal(fileRights)
	if err != nil{
		fmt.Println(err)
		return
	}
	jsonString := string(jsonBytes)
   	fmt.Println(jsonString)
	
	file_name := "D:\\Truffle_test\\user_files\\file_permissions_"+fileRights.FileId+".js" //最好在文件名加入就诊卡号以标识
	var file *os.File
	var flag bool
	_, err = os.Stat(file_name)
    if err == nil{
        flag = true
    }
	if os.IsNotExist(err){  //文件不存在
		flag = false
	}
	if flag == false {
		file, err = os.Create(file_name)
		if err != nil{
			fmt.Println("创建失败",err)
			return
		}		
	}else{
		file,err = os.OpenFile(file_name,os.O_RDWR,0666)
		if err!=nil{
			fmt.Println("打开文件错误：",err)
			return
		}
	}
	defer file.Close()
	n,err:=io.WriteString(file,jsonString)
	if err != nil {
		fmt.Println("写入错误：",err)
		return
	}
	fmt.Println("写入成功：n=",n)
	fmt.Fprintf(w, jsonString) 
}



func main() {
	token_map = make(map[string]string)  	
	http.HandleFunc("/getToken", obtainToken)
	//fmt.Println("Hello World!")
	http.HandleFunc("/writeEMC",sendEMR)
	http.HandleFunc("/writeEMCRights",sendEMRRights)
	log.Println("Cloud Server is working on port 8085")
	log.Fatal(http.ListenAndServe(":8085", nil)) 	
}
