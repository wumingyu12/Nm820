var page13_model=angular.module('MyApp.page13',[]);

//第三个表格24小时通风等级曲线的控制器，
page12_model.controller('page13_tongfenMaxminCtrl',[
	'$scope',
	'$http',
	function ($scope,$http){
		//用于切换是用图表显示还是表格
		$scope.ifShowTab=false;//一开始时是显示表格不显示曲线
		$scope.ifShowCurve=true;

		$http.get('/testjson/page13/tongfengClass.json').success(function(data){
			$scope.tongfengdata = data[0];
			console.log($scope.tongfengdata);
		});
		 //用于显示的数据，日龄，最小通风等级，最大通风等级
  		$scope.tongfengCurves=[
  			{day:'0',minClass:'12',maxClass:'13'},
  			{day:'0',minClass:'12',maxClass:'13'},
  			{day:'0',minClass:'12',maxClass:'13'},
  			{day:'0',minClass:'12',maxClass:'13'},
  			{day:'0',minClass:'12',maxClass:'13'},
  			{day:'0',minClass:'12',maxClass:'13'},
  			{day:'0',minClass:'12',maxClass:'13'},
  			{day:'0',minClass:'12',maxClass:'13'},
  			{day:'0',minClass:'12',maxClass:'13'},
  			{day:'0',minClass:'12',maxClass:'13'},
  			{day:'0',minClass:'12',maxClass:'13'},
  		]
  		//点击显示表格与图表切换
		$scope.changeTabChat = function(){
	    	$scope.ifShowTab=!$scope.ifShowTab;
	    	$scope.ifShowCurve=!$scope.ifShowCurve;
		};

	}
]);