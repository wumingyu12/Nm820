
var page12_model=angular.module('MyApp.page12', []);

/*===========================24小时温度曲线 =================================================
  //24小时温度曲线的控制器，
===================================================================================*/
page12_model.controller('page12_LineCtrl_wenduDay',[
  '$scope',
  '$http',
  function ($scope,$http){
    //图表点击时动作
    $scope.series=["室内温度","舍外温度"];
    $scope.data = [
      [10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10],
      [10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10],
      ];
    $scope.labels=[0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23];
    $http.get('/resetful/nm820Json/Get24TemHumi.json').success(function(data){
      var nowhour=data.Nowhour;//json返回的当前时间
      //以当前时间作为截取，让当前时间对应的温湿度值放到数组的最后，【0,1,2。。now| now+1，。。】变【now+1，。。。，前1小时，now】
      var nowfront=data.Tavg.slice(0,nowhour+1);
      var nowafter=data.Tavg.slice(nowhour+1);
      $scope.data[0]=nowafter.concat(nowfront);

      //时间进度条要变化
      var timefront=data.Time.slice(0,nowhour+1);
      var timeafter=data.Time.slice(nowhour+1);
      $scope.labels=timeafter.concat(timefront);
    });
  }
]);


/*====================月平均，最大，最小温度=================================
  参数：
      1.$scope.ocw用于初始化显示图表的数据
      2.$scope.isRotate=false;//一开始刷新图标是不转的
  方法：
      1.reflashData更新温度月数据
=========================================================*/
page12_model.controller('page12_LineCtrl_wenduMonth',[
  '$scope',
  '$http',
  function ($scope,$http){
    $scope.isRotate=false;//一开始刷新图标是不转的
      //图表点击时动作
    $scope.onClick = function (points, evt) {
        console.log(points, evt);
    };
    //图表初始化显示的东西
    $scope.ocw = {
      "series": ["最高温度","平均温度","最低温度"],
      //也可以用上面的scope.data的形式
      "data": [["20", "20", "20", "20", "20", "20", "20","20", "20", "20","20","20","20","20","20","20","20","20","20","20","20","20","20","20","20","20","20","20","20","20","20"],
          ["19", "19", "19", "19", "19", "19", "19","19", "19", "19","19","19","19","19","19","19","19","19","19","19","19","19","19","19","19","19","19","19","19","19","19"],
          ["18", "18", "18", "18", "18", "18", "18","18", "18", "18","18","18","18","18","18","18","18","18","18","18","18","18","18","18","18","18","18","18","18","18","18"]],
      "labels":["1", "2", "3", "4", "5", "6","7", "8", "9", "10", "11", "12", "13","14", "15", 
           "16", "17", "18", "19", "20","21","22","23","24","25","26","27","28","29","30","31"],
      "colours": [{ // default
        //填充颜色，有多个曲线时如果后面的没定义就随机
          "fillColor": ["rgba(22, 211, 112, 0)","rgba(22, 211, 112, 0)"],
          //图例颜色如seriesA：黄色
          "strokeColor": "rgba(20,100,13,1)",
          "pointColor": "rgba(220,220,220,1)",
          "pointStrokeColor": "#fff",
          "pointHighlightFill": "#fff",
          "pointHighlightStroke": "rgba(151,187,205,0.8)"
      }]
    };

    //resetful更新数据,单击刷新按钮时
    $scope.reflashData=function(){
      $scope.isRotate=true;//刷新图标是转的
      //让数据先初始化加载后有个动画效果
      $scope.ocw.data[0]=["20", "20", "20", "20", "20", "20", "20","20", "20", "20","20","20","20","20","20","20","20","20","20","20","20","20","20","20","20","20","20","20","20","20","20"];
      $scope.ocw.data[1]=["20", "20", "20", "20", "20", "20", "20","20", "20", "20","20","20","20","20","20","20","20","20","20","20","20","20","20","20","20","20","20","20","20","20","20"];
      $scope.ocw.data[2]=["20", "20", "20", "20", "20", "20", "20","20", "20", "20","20","20","20","20","20","20","20","20","20","20","20","20","20","20","20","20","20","20","20","20","20"];

      $http.get('/resetful/nm820/GetDataHistory/Tem').success(function(data){
        console.log(data.Avgs);
        $scope.ocw.labels=data.Days;
        $scope.ocw.data[0]=data.Maxs;
        $scope.ocw.data[1]=data.Avgs;
        $scope.ocw.data[2]=data.Mins;
        console.log($scope.ocw);
        $scope.isRotate=false;//刷新图标是不转的
      });
    }
  }
]);

/*===========================24小时湿度曲线 =================================================
  //24小时湿度曲线的控制器，
===================================================================================*/
page12_model.controller('page12_LineCtrl_shiduDay',[
  '$scope',
  '$http',
  function ($scope,$http){
    //图表点击时动作
    $scope.series=["室内湿度","舍外湿度"];
    $scope.data = [
      [70,70,70,70,70,70,70,70,70,70,70,70,70,70,70,70,70,70,70,70,70,70,70,70],
      [70,70,70,70,70,70,70,70,70,70,70,70,70,70,70,70,70,70,70,70,70,70,70,70],
      ];
    $scope.labels=[0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23];
    $http.get('/resetful/nm820Json/Get24TemHumi.json').success(function(data){
      var nowhour=data.Nowhour;//json返回的当前时间
      //以当前时间作为截取，让当前时间对应的温湿度值放到数组的最后，【0,1,2。。now| now+1，。。】变【now+1，。。。，前1小时，now】
      var nowfront=data.Havg.slice(0,nowhour+1);
      var nowafter=data.Havg.slice(nowhour+1);
      $scope.data[0]=nowafter.concat(nowfront);

      //时间进度条要变化
      var timefront=data.Time.slice(0,nowhour+1);
      var timeafter=data.Time.slice(nowhour+1);
      $scope.labels=timeafter.concat(timefront);
    });
  }
]);

/*====================月平均，最大，最小湿度=================================
  参数：
      1.$scope.ocw用于初始化显示图表的数据
      2.$scope.isRotate=false;//一开始刷新图标是不转的
  方法：
      1.reflashData更新温度月数据
=========================================================*/
page12_model.controller('page12_LineCtrl_shiduMonth',[
	'$scope',
  '$http',
	function ($scope,$http){
	  	//图表点击时动作
		$scope.onClick = function (points, evt) {
	    	console.log(points, evt);
		};
		$scope.data = [
			[80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80],
			[80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80],
      [80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80],
	  	];
	  	$scope.series=["最高湿度","平均湿度","最低湿度"];
	  	$scope.labels=[1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30];

          //resetful更新数据,单击刷新按钮时
    $scope.reflashData=function(){
      $scope.isRotate=true;//刷新图标是转的
      //让数据先初始化加载后有个动画效果
      $scope.data[0]=[80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80];
      $scope.data[1]=[80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80];
      $scope.data[2]=[80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80,80];

      $http.get('/resetful/nm820/GetDataHistory/Humi').success(function(data){
        console.log(data.Avgs);
        $scope.labels=data.Days;
        $scope.data[0]=data.Maxs;
        $scope.data[1]=data.Avgs;
        $scope.data[2]=data.Mins;
        $scope.isRotate=false;//刷新图标是不转的
      });
    }
	}
]);

/*====================系统参数总集成，请注意这里的cotrol将会被多个表使用=================================
  参数：
    1.relayList 用来继电器编码的下拉选项
  方法：
      1.reflashData更新温度月数据
=========================================================*/
page12_model.controller('page12_SysvalCtrl', [
  '$scope',
  '$http',
  '$filter',
  function ($scope,$http,$filter){
    $http.get('/resetful/nm820/sysPara/SysValTable').success(function(data){
      $scope.sysVal=data;
      console.log($scope.sysVal);
    });

    //===============下面的代码在继电器编码的选项中用到=======================
    $scope.relayList = [
      {value: 0, text: '0--手动'},
      {value: 1, text: '1--风机组1'},
      {value: 2, text: '2--风机组2'},
      {value: 3, text: '3--风机组3'},
      {value: 4, text: '4--风机组4'},
      {value: 5, text: '5--风机组5'},
      {value: 6, text: '6--风机组6'},
      {value: 7, text: '7--风机组7'},
      {value: 8, text: '8--风机组8'},
      {value: 9, text: '9--加热'},
      {value: 10, text: '10--冷却水泵'},
      {value: 11, text: '11--喷雾'},
      {value: 12, text: '12--回流阀'},
      {value: 13, text: '13--照明1'},
      {value: 14, text: '14--照明2'},
      {value: 15, text: '15--侧风窗开'},
      {value: 16, text: '16--侧风窗关'},
      {value: 17, text: '17--幕帘开'},
      {value: 18, text: '18--幕帘关'},
      {value: 19, text: '19--卷帘1开'},
      {value: 20, text: '20--卷帘1关'},
      {value: 21, text: '21--卷帘2开'},
      {value: 22, text: '22--卷帘2关'},
      {value: 23, text: '23--卷帘3开'},
      {value: 24, text: '24--卷帘3关'},
      {value: 25, text: '25--卷帘4开'},
      {value: 26, text: '26--卷帘4关'},
      {value: 27, text: '27--喂料'},
      {value: 28, text: '28--额外系统1'},
      {value: 29, text: '29--额外系统2'},
      {value: 30, text: '30--额外系统3'},
      {value: 31, text: '31--额外系统4'}
    ]; 

    $scope.relayShow = function(myvalue) {
      var selected = $filter('filter')($scope.relayList, {value: myvalue});//过滤器在relaylist中找出value=myvalue的对象
      return (myvalue && selected.length) ? selected[0].text : 'Not set';
    };
    //====================================================================
  }
]);

