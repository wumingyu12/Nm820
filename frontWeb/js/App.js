var myApp = angular.module('MyApp', ['ngAnimate','MyApp.page1','MyApp.page2']);//[]里可以注入模块

myApp.controller('bodyCtrl',[
	'$scope',
	'$rootScope',
	function ($scope,$rootScope){
		$rootScope.mainShows=[//主页面用
			{url:'page1_state/page1_controlIndex.html'},
			{url:'page2_control/index.html'}
		];
		$rootScope.mainShow=$rootScope.mainShows[0];
		//================================全局变量==============================
		//设备列表，在page2 中会用到
		$rootScope.equiments=[//设备集合
			{eqname:"风机1",icon:"icon-wendu"},
			{eqname:"风机1",icon:"icon-wendu"}, 
			{eqname:"风机1",icon:"icon-wendu"}, 
			{eqname:"风机1",icon:"icon-wendu"}, 
			{eqname:"风机1",icon:"icon-wendu"}, 
			{eqname:"风机1",icon:"icon-wendu"},
			{eqname:"风机1",icon:"icon-wendu"}, 
			{eqname:"风机1",icon:"icon-wendu"}, 
			{eqname:"风机1",icon:"icon-wendu"}, 
			{eqname:"风机1",icon:"icon-wendu"}, 
			
		];
	}
]);


myApp.controller('navleftCtrl',[
	'$scope',
	'$rootScope',
	function ($scope,$rootScope){
		$scope.currentIndexNum=0;//指示当前选择的导航按钮，用于ng-class高亮显示
		$scope.showMain=function(index){//单击导航条的按钮显示对应页面
			$rootScope.mainShow=$rootScope.mainShows[index];
			$scope.navShow=!$scope.navShow;//隐藏回导航条
			$scope.currentIndexNum=index;
			console.log($scope.currentIndexNum);
		};
	}
]);



