var page13_model=angular.module('MyApp.page13',[]);

/*============================控制器 温度曲线======================================
	1.图表切换
	2.请求温度曲线数据
=======================================================================================*/
page13_model.controller('page13_wenduCurveCtrl', [
	'$scope',
	'$http',
	function ($scope,$http){
		$scope.ifShowTab=true;//一开始时是显示表格不显示曲线
		$scope.ifShowCurve=false;

		//图表显示数据
		$scope.CurveData={
			"series": ["目标温度","加热温度","制冷温度"],
        	//也可以用上面的scope.data的形式
        	"data": [[1,2,3,45,6,7,7,8,8],
        			 [2,34,5,5,5,5,5,5,5],
        			 [3,4,5,5,66,77,88,77]],
        	"labels":[1,2,3,4,5,6,7,8,9],
		};

		$http.get('/resetful/nm820/sysPara/WenduCurve').success(function(data){
			$scope.Days=data.Day;
			$scope.Targets=data.Target;
			$scope.Heats=data.Heat;
			$scope.Cools=data.Cool;

			//温度曲线图表的样子
			$scope.CurveData.labels=data.Day;
			$scope.CurveData.data[0]=data.Target;
			$scope.CurveData.data[1]=data.Heat;
			$scope.CurveData.data[2]=data.Cool;
		});


		//点击显示表格与图表切换
		$scope.changeTabChat = function(){
	    	$scope.ifShowTab=!$scope.ifShowTab;
	    	$scope.ifShowCurve=!$scope.ifShowCurve;
		};


	}
])
//第三个表格24小时通风等级曲线的控制器，
page12_model.controller('page13_tongfenMaxminCtrl',[
	'$scope',
	'$http',
	function ($scope,$http){
		//用于切换是用图表显示还是表格


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