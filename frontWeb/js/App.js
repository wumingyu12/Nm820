var myApp = angular.module('MyApp', ['ngAnimate','MyApp.page1','MyApp.page2','chart.js']);//[]里可以注入模块

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

myApp.controller('LineCtrl',[
	'$scope',
	function ($scope){
		$scope.labels = ["January", "February", "March", "April", "May", "June", "July"];
		$scope.series = ['Series A', 'Series B'];
		$scope.data = [
			[65, 59, 80, 81, 56, 55, 40],
	    	[28, 48, 40, 19, 86, 27, 90]
	  	];
		$scope.onClick = function (points, evt) {
	    	console.log(points, evt);
		};
	}
]);




