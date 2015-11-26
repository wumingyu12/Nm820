var myApp = angular.module('MyApp', [
	'ui.router',
	'ngAnimate',
	'MyApp.page1',
	'MyApp.page12',
	'MyApp.page11',
	'MyApp.page13',
	'MyApp.page51',
	'MyApp.page2',
	'chart.js'
]);//[]里可以注入模块

//================================
//ui-view的路由
//==============================
myApp.config(function($stateProvider, $urlRouterProvider) {
	$urlRouterProvider.otherwise("/page1");//默认
	$stateProvider
		//========================页面1==========================
        .state('page1', {//环境控制监控顶级页面
            url: "/page1",
            templateUrl: "page1_state/page1_controlIndex.html"
        })
            .state('page1.contState', {//环境控制监控,第二级页面状态监控
	            url: "/contState",
	            templateUrl: "page1_state/page11_contState.html"
        	})
        	.state('page1.historyData', {//环境控制监控,第二级页面历史数据
	            url: "/historyData",
	            templateUrl: "page1_state/page12_historyDate.html"
        	})
        	.state('page1.parmShow', {//环境控制监控,第二级页面参数设置
	            url: "/parmShow",
	            templateUrl: "page1_state/page13_parmShow.html"
        	})

        //================页面二================================
        .state('page2', {
            url: "/page2",
            templateUrl: "page1_state/page1_controlIndex.html"
        })
        .state('page3', {
            url: "/page3",
            templateUrl: "page1_state/page1_controlIndex.html"
        })
        .state('page4', {
            url: "/page4",
            templateUrl: "page2_control/index.html"
        })
        .state('page5', {//视频监控顶级页面
            url: "/page5",
            templateUrl: "page5_video/page51_videoShow.html"
        });;
})

myApp.controller('bodyCtrl',[
	'$scope',
	'$rootScope',
	function ($scope,$rootScope){
		$rootScope.mainShows=[//主页面用
			{url:'page1_state/page1_controlIndex.html'},
			{url:'page2_control/index.html'},
			{url:'page2_control/index.html'},
			{url:'page2_control/index.html'},
			{url:'page5_video/page51_videoShow.html'}
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






