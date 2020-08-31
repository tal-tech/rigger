package internal

import (
	"bytes"
	"errors"
	"fmt"
	"go/parser"
	"go/token"

	"github.com/tal-tech/rigger/common"
)

var (
	phpCodeTpl = `<?php
//usage $url gateway addr
//$client = new %sClient($url="http://127.0.0.1:8080/");

//$resp = $client->SayHello('php client');

//if ($resp->getHttpStatus() != '200') {
//	echo 'http error';
//} else if($resp->hasError()) {
//	echo $resp->getErrorMessage();
//} else {
//	echo $resp->getData();
//}

class Response {
	private $_hasError = false;
	private $_errmsg;
	private $_body;
	private $_status;
	private $_protocal;
	private $_messageId;   

	public function __construct($response, $mid) {
		$this->parse($response);
		$this->_messageId = $mid;
	}

	private function parse($response) {
		list($header, $body) = explode("\r\n\r\n", $response, 2);
		
		$headerArray = explode("\r\n", $header);
		
		list($protocal,$status, $message) = explode(" ", $headerArray[0]);
		
		$this->_status = $status;
		$this->_protocal = $protocal;
		
		unset($headerArray[0]);

		foreach ($headerArray as $value) {
			list($headerKey, $headerValue) = explode(":", $value, 2);
			$this->_header[trim($headerKey)] = trim($headerValue);
		}

		if ($this->_header["X-Rpcx-Messagestatustype"] == "Error") {
			$this->_hasError = true;
			$this->_errmsg = $this->_header["X-Rpcx-Errormessage"];
		}
		
		$this->_body = $body;
	}

	public function hasError() {
		return $this->_hasError;
	}

	public function getData() {
		return $this->_body;
	}

	public function getErrorMessage() {
		return $this->_errmsg;
	}

	public function getHttpStatus() {
		return $this->_status;
	}
}

class %sClient {
	
	const ServiceName = "%s";

	private $_url;
	private $_httpClient;
	private $_timeout;
	
	public function __construct($url, $timeout = 3) {
		$this->_url = $url;
		$this->_timeout = $timeout;
	}

	private function call($uri, $data) {
		$curl = curl_init();
		$msgId = $this->getUuid();
		
		$headers = [	
			"Content-type:application/json",
		];

		$curlOpt = [
			CURLOPT_URL => $this->_url . $uri,
			CURLOPT_USERAGENT => 'Xes-Micro HttpClient/Curl',
			CURLOPT_CONNECTTIMEOUT => $this->_timeout,
			CURLOPT_POST => true,
			CURLOPT_RETURNTRANSFER => 1,
			CURLOPT_HTTPHEADER => $headers,
			CURLOPT_POSTFIELDS => json_encode($data),
			CURLOPT_HEADER => true,
			CURLOPT_NOBODY => false,
		];
		
		curl_setopt_array($curl, $curlOpt);

		$response = curl_exec($curl);

		curl_close($curl);
   
		return (new Response($response, $msgId));
	}
`
)

func GenPHPHttpClient(serviceFile string) (*bytes.Buffer, error) {
	exists, err := common.PathExists(serviceFile)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("file " + serviceFile + " not exist")
	}

	fset := token.NewFileSet()

	fs, err := parser.ParseFile(fset, serviceFile, nil, parser.ParseComments)

	if err != nil {
		return nil, err
	}

	serviceName := getServiceName(fs)

	buffer := bytes.NewBufferString(genPhpClass(serviceName) + "\n")

	for _, decl := range fs.Decls {
		fn := parseFunc(decl)

		if fn == nil {
			continue

		}
		if len(fn.Args) != 3 {
			continue
		}

		buffer.WriteString(genPhpFunc(serviceName, fn.Name) + "\n")
	}

	buffer.WriteString("}\n")

	return buffer, nil
}

func genPhpFunc(serviceName string, fn string) string {
	body := fmt.Sprintf(`
	// todo 将arg替换为真实的输入参数
	public function %s($arg) {
		$data = [
			'arg' => $arg, 
		];
		return $this->call("/%s/%s", $data);
	}`, fn, serviceName, fn)

	return body
}

func genPhpClass(serviceName string) string {
	return fmt.Sprintf(phpCodeTpl, serviceName, serviceName, serviceName)
}
