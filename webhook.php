<?php

print_r($_REQUEST);
$key =$_REQUEST["data"]["USER_ID"];
$url="https://10.1.41.20/rest/1/ip0yvzl6zxbbqc2s/user.get?id=".$key;   //url to bitrix function api
$ch = curl_init();
curl_setopt($ch, CURLOPT_URL, $url);
curl_setopt($ch, CURLOPT_RETURNTRANSFER, 1);
curl_setopt($ch, CURLOPT_SSL_VERIFYHOST, 0);
curl_setopt($ch, CURLOPT_SSL_VERIFYPEER, 0);
$result = curl_exec($ch);
$vars = json_decode($result, true);
$exten=$vars["result"][0]["UF_PHONE_INNER"];
$_REQUEST["data"]["EXTENSION"]=$exten; 
//=================================== Redirect to asterlink after fix==============
$url="http://localhost:801/originate/"; //Path to Asterlink
//$url="http://41.234.8.238:5678/originate/";
$postdata = http_build_query($_REQUEST);

$ch = curl_init();
curl_setopt($ch, CURLOPT_URL, $url);
curl_setopt($ch, CURLOPT_POST, true);
curl_setopt($ch, CURLOPT_POSTFIELDS, $postdata);
curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
$result = curl_exec($ch);

//==================================================================================

