//page14 视频显示与ptz控制
var page51_model=angular.module('MyApp.page51', []);

//视频1的，ptz控制等
page51_model.controller('page51_mainvideoCtrl', [
	'$scope',
	'$http', 
	function ($scope,$http){
		//ptz 控制函数
		$scope.hkptzResetful=function(camnum,mode,speed){//resetful camnum代表摄像头号数，mode代表上下左右远近停止的调节，speed代表速度
			$http.get("/resetful/hkPtz/Continuous/"+ camnum +"/"+ mode +"/"+speed) //60代表输出，1代表1号
		}
	}
]);