
var page12_model=angular.module('MyApp.page12', []);

page12_model.controller('page12_LineCtrl',[
	'$scope',
	function ($scope){
		$scope.labels = ["0时", "1时", "2时", "3时", "4时", "5时", "6时","7时", "8时", "9时", "10时", "11时", "12时", "13时","14时", "15时", "16时", "17时", "18时", "19时", "20时","21时","22时","23时"];
		$scope.series = ['温度', '湿度'];
		$scope.options={
			scaleLineColor : "rgba(33,0,111,.1)",
		};
		$scope.data = [
			[65, 59, 80, 81, 56, 55, 40,65, 59, 80, 81, 56, 55, 40,65, 59, 80, 81, 56, 55, 40,22,23,24],
	    	[22, 59, 80, 33, 56, 55, 40,22, 59, 80, 81, 55, 55, 40,66, 59, 80, 81, 33, 55, 40,22,23,24],
	  	];
		$scope.onClick = function (points, evt) {
	    	console.log(points, evt);
		};
		$scope.cc=[{
         "fillColor": "rgba(224, 108, 112, 1)",
         "strokeColor": "rgba(207,100,103,1)",
         "pointColor": "rgba(220,220,220,1)",
         "pointStrokeColor": "#fff",
         "pointHighlightFill": "#fff",
         "pointHighlightStroke": "rgba(151,187,205,0.8)"
       }];

        var json = {
        	"series": ["SeriesA"],
        	"data": [["90", "99", "80", "91", "76", "75", "60", "67", "59", "55"]],
        	"labels": ["01", "02", "03", "04", "05", "06", "07", "08", "09", "10"],
        	"colours": [{ // default
          		"fillColor": "rgba(224, 108, 112, 1)",
          		"strokeColor": "rgba(207,100,103,1)",
          		"pointColor": "rgba(220,220,220,1)",
          		"pointStrokeColor": "#fff",
          		"pointHighlightFill": "#fff",
          		"pointHighlightStroke": "rgba(151,187,205,0.8)"
        	}]
      	};
  		$scope.ocw = json;
	}
]);