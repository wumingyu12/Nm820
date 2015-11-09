//page14 视频显示与ptz控制
var page14_model=angular.module('MyApp.page14', []);

//视频1的，ptz控制等
page14_model.controller('page14_video1Ctrl', [
	'$scope',
	'$http', 
	function ($scope,$http){
		//ptz 控制函数
		$scope.ptzUp=function(){
			$http.put("")
		}
	}
]);