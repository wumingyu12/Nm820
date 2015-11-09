var myApp = angular.module('MyApp', [
	'ngAnimate',
	'MyApp.page1',
	'MyApp.page12',
	'MyApp.page11',
	'MyApp.page13',
	'MyApp.page14',
	'MyApp.page2',
	'chart.js'
]);//[]里可以注入模块

myApp.controller('bodyCtrl',[
	'$scope',
	'$rootScope',
	function ($scope,$rootScope){
		$rootScope.mainShows=[//主页面用
			{url:'page1_state/page1_controlIndex.html'},
			{url:'page2_control/index.html'}
		];
		$rootScope.mainShow=$rootScope.mainShows[0];
	}
]);


myApp.controller('navleftCtrl',[
	'$scope',
	'$rootScope',
	function ($scope,$rootScope){
		$scope.currentIndexNum=0;//指示当前选择的导航按钮，用于ng-class高亮显示
		$scope.showMain=function(index){//单击导航条的按钮显示对应页面
			$rootScope.mainShow=$rootScope.mainShows[index];
			$scope.currentIndexNum=index;//将当前高亮显示的index改变
		};
	}
]);






