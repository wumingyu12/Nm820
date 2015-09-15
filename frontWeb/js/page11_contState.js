var page11_model=angular.module('MyApp.page11', []);


/*========================服务 在非当前页面中，停止心跳获取Nm820状态===============================
  动作：
    1.当在本页面那就定时发送请求resetful
    2.当当前界面不是本界面就不发送心跳包
  注入的地方：
    1.page11StatepanCtrl 状态显示控制器
    2.
  提供的变量：
    1.nm820StateDate={TemAvg:"0",HumiAvg:"0",GDay:"0",FanLevel:"0",Year:"0",Month:"0",Day:"0",Hour:"0",Min:"0",Sec:"0"};
     TemAvg 温度  HumiAvg 湿度  GDay 日龄 FanLevel 通风等级
     Year Month Day Hour Min Sec 时间
==================================================================================================*/
page11_model.factory('page11getstateSer', ['$timeout','$http',function($timeout,$http){
  var addzero=function(s){ //时间补零函数
    return s < 10 ? '0' + s: s;
  };
  //保持状态信息的类
  var nm820StateDate={TemAvg:"0",HumiAvg:"0",GDay:"0",FanLevel:"0",Year:"0",Month:"0",Day:"0",Hour:"0",Min:"0",Sec:"0"};
  //定时更新数据，通过resetful
  var longPoll = function() {
    console.log("心跳");
    $http.get("/resetful/nm820/GetState")
      .success(function(data) { 
        nm820StateDate.TemAvg=data.TemAvg/10;//返回的数值是乘上10的
        nm820StateDate.HumiAvg=data.HumiAvg/10;
        nm820StateDate.GDay=data.GDay; //日龄
        nm820StateDate.FanLevel=data.FanLevel; //通风等级
        nm820StateDate.Year=data.Year;
        nm820StateDate.Month=addzero(data.Month);
        nm820StateDate.Day=addzero(data.Day);
        nm820StateDate.Hour=addzero(data.Hour);
        nm820StateDate.Min=addzero(data.Min);
        nm820StateDate.Sec=addzero(data.Sec);
        console.log(nm820StateDate);
      })
      .error(function(){
        nm820StateDate.TemAvg="fff";
      });
    $timeout(longPoll, 1000);
  }; 
  longPoll();//记得一开始要启动定时
  return {
    nm820StateDate:nm820StateDate,
  };
}]);

/*=====================指令 缩放并设置定位位置================================================
 依赖：1.jquery库
       2.一个父节点，并且子节点是相对于父节点绝对定位(在css里面有position: absolute;)
 作用：1.当原图长宽为parentWidth与parentHeight是元素的top和left，
       按照父元素的长宽按比例缩放videoWidth，videoHeight。和摆放位置
     2.改变值：top left width height
       
  参数：
    posTop:'@',//无缩放时的top
    posLeft:'@',//无缩放时的left
    parentWidth:'@',//无缩放时父节点的长
    parentHeight:'@',//无缩放时父节点的高
    myWidth:'@',//无缩放时我的长
    myHeight:'@',//无缩放我的高
 使用示范：
        set-position 
        pos-top="236" 
        pos-left="3279" 
        parent-width="4096" 
        parent-height="2304" 
        my-widht="110" 
        my-height="110"
  遗留问题：
      1.video缩放时的width与height是等比例的，意味着不能拉伸 
==============================================================================================*/

page11_model.directive('page11SetPosition',function(){
  // Runs during compile
  return {
    // name: '',
    // priority: 1,
    // terminal: true,
    //
    scope: {
      posTop:'@',//无缩放时的top
      posLeft:'@',//无缩放时的left
      parentWidth:'@',//无缩放时父节点的长
      parentHeight:'@',//无缩放时父节点的高
      myWidth:'@',//无缩放时我的长
      myHeight:'@',//无缩放我的高
    }, // {} = isolate, true = child, false/undefined = no change
    // controller: function($scope, $element, $attrs, $transclude) {},
    // require: 'ngModel', // Array = multiple requires, ? = optional, ^ = check parent elements
     restrict: 'EA', // E = Element, A = Attribute, C = Class, M = Comment
    // template: '',
    // templateUrl: '',
    // replace: true,
    // transclude: true,
    // compile: function(tElement, tAttrs, function transclude(function(scope, cloneLinkingFn){ return function linking(scope, elm, attrs){}})),
    link: function(scope, elem, attrs, controller) {
      //计算缩放比例
      //console.log(elem.parent());
      var scaleWdith=parseFloat(elem.parent().css('width'))/parseFloat(scope.parentWidth);
      var scaleHeight=parseFloat(elem.parent().css('height'))/parseFloat(scope.parentHeight);
      //alert(scaleHeight+"---"+scaleWdith)
      //缩放本元素大小
      var currentWidth=parseFloat(scaleWdith*parseFloat(scope.myWidth));
      var currentHeight=parseFloat(scaleHeight*parseFloat(scope.myHeight));
      elem.css('width',currentWidth+'px');
      elem.css('height',currentHeight+'px');
      //修改缩放后的位置
      var currentLeft=parseFloat(scaleWdith*parseFloat(scope.posLeft));
      var currentTop=parseFloat(scaleHeight*parseFloat(scope.posTop));
      elem.css('left',currentLeft+'px');
      elem.css('top',currentTop+'px');
      //scope.$apply();//没了这个是不会在template里面更新的。
    }
  };
});

/*
//24小时温度曲线的控制器，
page11_model.controller('page11_LineCtrl_wenduDay',[
	'$scope',
	function ($scope){
		//如果不用json图表的另一种表达数据的方式
		//$scope.data = [
		//	[65, 59, 80, 81, 56, 55, 40,65, 59, 80, 81, 56, 55, 40,65, 59, 80, 81, 56, 55, 40,22,23,24],
	  //  	[22, 59, 80, 33, 56, 55, 40,22, 59, 80, 81, 55, 55, 40,66, 59, 80, 81, 33, 55, 40,22,23,24],
	  //	];
	
	  	//图表点击时动作
		$scope.onClick = function (points, evt) {
	    	console.log(points, evt);
		};


    var lineJson = {
    	"series": ["室内温度"],
    	//也可以用上面的scope.data的形式
    	"data": [["18", "19", "20", "21", "21", "22", "21.4","22.4", "24.5", "22.7","23.8","25.7","22.4","26.7","22.9","24.5","19.5","26","21","23","22","22","30","15"]],
    	"labels":["", "", "", "", "", "", "","", "", "", "", "", "", "","", "", "", "", "", "", "","","",""],
    	//"colours": [{ // default,可以在canvas里面通过colours="ocw.colours"-使用
    		//填充颜色，有多个曲线时如果后面的没定义就随机
      		//"fillColor": ["rgba(22, 211, 112, 1)"],
      		//图例颜色如seriesA：黄色
      		//"strokeColor": "rgba(20,100,13,1)",
      		//"pointColor": "rgba(220,220,220,1)",
      		//"pointStrokeColor": "#fff",
      		//"pointHighlightFill": "#fff",
      		//"pointHighlightStroke": "rgba(151,187,205,0.8)"
    	//}]
    };
    //动态刷新图表
    $scope.wenduCurrent="0";
		$scope.ocw = lineJson;
    var testtimeout=function(){
      var tem=lineJson.data[0][0];
      $scope.wenduCurrent=tem;
      lineJson.data[0].shift();//最前面的数移走，返回新的数列
      lineJson.data[0].push(tem);
      console.log(lineJson.data[0])
      delete tem;//用完后记得释放，否则每一定时都会生成一个tem
    };
    setInterval(function(){
      $scope.$apply(testtimeout);
    },1000);
    testtimeout();
	}
]);
*/

/*========================控制器==温湿度状态显示面板==============================================
  双向绑定：1.TemAvg 温度  HumiAvg 湿度  GDay 日龄 FanLevel 通风等级
            2. Year Month Day Hour Min Sec 时间,通过服务
  内部方法：
=================================================================================================*/
page11_model.controller('page11StatepanCtrl', [
  '$scope',
  '$http',
  '$timeout',
  'page11getstateSer',
  function($scope,$http,$timeout,page11getstateSer){
    //定时更新数据，通过从服务page11getstateSer中
    //用watch更新数据
    $scope.$watch(function () { return page11getstateSer.nm820StateDate; },
      function (value) {
          //console.log("In $watch - lastUpdated:" + value);
          $scope.nm820state = value;
      }
    );


  }
])

