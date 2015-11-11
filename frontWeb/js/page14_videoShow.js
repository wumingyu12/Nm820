//page14 视频显示与ptz控制
var page14_model=angular.module('MyApp.page14', []);

//视频1的，ptz控制等
page14_model.controller('page14_mainvideoCtrl', [
	'$scope',
	'$http', 
	function ($scope,$http){
		//ptz 控制函数
		$scope.hkptzResetful=function(camnum,mode,speed){//resetful camnum代表摄像头号数，mode代表上下左右远近停止的调节，speed代表速度
			$http.get("/resetful/hkPtz/Continuous/"+ camnum +"/"+ mode +"/"+speed) //60代表输出，1代表1号
		}
	}
]);