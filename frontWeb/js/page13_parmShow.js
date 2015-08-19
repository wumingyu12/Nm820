var page13_model=angular.module('MyApp.page13',[]);

//24小时温度曲线的控制器，
page12_model.controller('page13_testCtrl',[
	'$scope',
	function ($scope){
		$scope.onClick = function(){
	    	console.log("points, evt");
		};


        var lineJson = {
        	"series": ["室内温度","室外温度"],
        	//也可以用上面的scope.data的形式
        	"data": [["18", "19", "20", "21", "21", "22", "21.4","22.4", "24.5", "22.7","23.8","25.7","22.4","26.7","22.9","24.5","19.5","26","21","23","22","22","30","15"],
        			["22.4","26.7","22.9","24.5","19.5","26","21","23","22","22","30","15","18", "19", "20", "21", "21", "22", "21.4","22.4", "24.5", "22.7","23.8","25.7"]],
        	"labels":["0时", "1时", "2时", "3时", "4时", "5时", "6时","7时", "8时", "9时", "10时", "11时", "12时", "13时","14时", "15时", "16时", "17时", "18时", "19时", "20时","21时","22时","23时"],
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
  		$scope.ocw = lineJson;
	}
]);