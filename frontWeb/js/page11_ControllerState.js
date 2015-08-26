var page11_model=angular.module('MyApp.page11', []);

//24小时温度曲线的控制器，
page11_model.controller('page11_LineCtrl_wenduDay',[
	'$scope',
	function ($scope){
		/*如果不用json图表的另一种表达数据的方式
		$scope.data = [
			[65, 59, 80, 81, 56, 55, 40,65, 59, 80, 81, 56, 55, 40,65, 59, 80, 81, 56, 55, 40,22,23,24],
	    	[22, 59, 80, 33, 56, 55, 40,22, 59, 80, 81, 55, 55, 40,66, 59, 80, 81, 33, 55, 40,22,23,24],
	  	];
	  	*/
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