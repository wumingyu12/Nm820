var page13_model=angular.module('MyApp.page13',["xeditable"]);

page13_model.run(function(editableOptions) {
  editableOptions.theme = 'bs3'; // bootstrap3 theme. Can be also 'bs2', 'default'
});

/*==================================指令 让表格tr可以编辑=========================================
	作用：用于编辑参数表
		1.双击会进入编辑模式
		2.悬停会改变背景颜色
	依赖：1.jquery库
		  2.
================================================================================================*/
page13_model.directive('page13EditGrid', function(){
	// Runs during compile
	return {
		// name: '',
		// priority: 1,
		// terminal: true,
		// scope: {}, // {} = isolate, true = child, false/undefined = no change
		// controller: function($scope, $element, $attrs, $transclude) {},
		// require: 'ngModel', // Array = multiple requires, ? = optional, ^ = check parent elements
		restrict: 'EA', // E = Element, A = Attribute, C = Class, M = Comment
		template: '<input type="text" class="input-medium" size="5">',
		// templateUrl: '',
		// replace: true,
		// transclude: true,
		// compile: function(tElement, tAttrs, function transclude(function(scope, cloneLinkingFn){ return function linking(scope, elm, attrs){}})),
		link: function(scope, elem, attrs, controller) {
			elem.bind("click",function(){
				alert("11111111");
			});

			//先保存原来的背景颜色
			var origenback=null;
			//悬停改变背景颜色
			elem.hover(	
				//进入函数
				function(){
					origenback=elem.css("background-color");//先保存原来的背景色
					elem.css("background-color","#cceeff");
				},
				//退出函数
				function(){
					elem.css("background-color",origenback);
				}
			);
		}
	};
});

/*================================通风等级表============================================
==============================================================================*/
page13_model.controller('page13_wenduTableCtrl', [
	'$scope',
	'$http',
	function ($scope,$http){
		$http.get('/resetful/nm820/sysPara/WindTables').success(function(data){
			$scope.data=data.WindTables;
			console.log($scope.data);
		});
	}
]);
/*===================================过滤器 通风等级表的fan属性的显示用========================================================
		将 int转换为byte
		如1 转换为 10000000
		  3 转换为 12000000
		  127      12345670
		  255      12345678
===============================================================================================*/
page13_model.filter('page13_FanIntToString_Filt',function(){
	return function(input){
		var str="";
		if((input&1)==1){//如果低位有值 如100111&000001 =0000001
			str=str+'1';
		}else{
			str=str+'0';
		};

		if((input&2)==2){//如果第二位有值 
			str=str+'2';
		}else{
			str=str+'0';
		};

		if((input&4)==4){//第三位
			str=str+'3';
		}else{
			str=str+'0';
		};

		if((input&8)==8){//第三位
			str=str+'4';
		}else{
			str=str+'0';
		};


		if((input&16)==16){//第三位
			str=str+'5';
		}else{
			str=str+'0';
		};


		if((input&32)==32){//第三位
			str=str+'6';
		}else{
			str=str+'0';
		};


		if((input&64)==64){//第三位
			str=str+'7';
		}else{
			str=str+'0';
		};


		if((input&128)==128){//第三位
			str=str+'8';
		}else{
			str=str+'0';
		};
		return str;
	};
});
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
]);
/*===========================控制器 最大最小通风等级===============================================
	第三个表格通风等级曲线的控制器，
============================================================================*/
page13_model.controller('page13_WindLevelCtrl', [
	'$scope',
	'$http',
	function ($scope,$http){
		$scope.ifShowTab=true;//一开始时是显示表格不显示曲线
		$scope.ifShowCurve=false;

		//图表显示数据
		$scope.CurveData={
			"series": ["最小通风等级","最大通风等级"],
        	//也可以用上面的scope.data的形式
        	"data": [[1,2,3,45,6,7,7,8,8],
        			 [2,34,5,5,5,5,5,5,5]],
        	"labels":[1,2,3,4,5,6,7,8,9],
		};

		$http.get('/resetful/nm820/sysPara/WindLevel').success(function(data){
			$scope.Day=data.Day;
			$scope.Min=data.Min;
			$scope.Max=data.Max;

			//温度曲线图表的样子
			$scope.CurveData.labels=data.Day;
			$scope.CurveData.data[0]=data.Min;
			$scope.CurveData.data[1]=data.Max;
		});


		//点击显示表格与图表切换
		$scope.changeTabChat = function(){
	    	$scope.ifShowTab=!$scope.ifShowTab;
	    	$scope.ifShowCurve=!$scope.ifShowCurve;
		};

		$scope.test=function(){
			 return "Username should be `awesome` or `error`";
		};
	}
]);