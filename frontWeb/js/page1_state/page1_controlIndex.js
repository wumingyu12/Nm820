var page1_model=angular.module('MyApp.page1', []);
page1_model.controller('page1_MainCtrl',[
	'$scope',
	'page11getstateSer',
	function ($scope,page11getstateSer){
		$scope.mainShows=[//主页面用
			{url:'page1_state/page11_contState.html'},
			{url:'page1_state/page12_historyDate.html'},
			{url:'page1_state/page13_parmShow.html'},
		];
		$scope.mainShow=$scope.mainShows[0];//默认是第一个页面

		$scope.currentIndexNum=0;//指示当前选择的导航按钮，用于ng-class高亮显示
		$scope.showMain=function(index){//单击导航条的按钮显示对应页面
			$scope.mainShow=$scope.mainShows[index];
			$scope.currentIndexNum=index;//将当前高亮显示的index改变

			//如果是状态显示页面，那就要不断resetful状态值
			if ($scope.currentIndexNum==0) {
				page11getstateSer.ifCurrentPage="yes";
			}else{
				page11getstateSer.ifCurrentPage="no";
			};
		};
	}
]);
