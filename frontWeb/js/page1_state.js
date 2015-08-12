
var model=angular.module('MyApp.page1', []);
/*指令圆环指示条
	CSS3圆环倒计时效果 百度
	http://www.yangqq.com/web/css3demo/index.html
	1用"border-radius"实现圆形
	2用"clip:rect"遮罩为半圆
	3 父级增加一层DIV
	4把父级DIV用"clip:rect"遮罩为一半
	5用"transform:rotate"连续改变扇形旋转角度
	6同理复制左半边扇形旋转
	7，0到180度右半圆转入左半边，左半圆不动
		180度到360度，右半圆完全转入左半边的情况下完全不动，左半圆转入右半边
		360度后回到第一步

	使用示范
	<ringstate stvalue='stvalue' 
			min-stvalue='10' 
			max-stvalue='80' 
			iconclass='icon-wendu'
			type='温度'
			unit='℃'
			></ringstate>
*/
model.directive('ringstate',function(){	
	return{
		restrict:'EA',
		scope:{
			stvalue:'=stvalue',//状态条要显示的值
			minStvalue:'@',//状态条的最小值和最大值可以确定stvalue当前的百分比
			maxStvalue:'@',
			//iconclass:'@',//图标的class，需要导入iconmoon的css
			//type:'@',//类型，text显示的名字,如温度，湿度
			//unit:'@',//如ppm ℃等
		},
		template:
			//底板，百分百布满父div，用来确定长宽的最少值，来生成正方形
			'<div style="height:100%;width:100%;">'+
			//为了半圆裁剪，加position:relative,底圆限制为正方形，这里的颜色将会percent比,颜色由init函数根据百分比确定
			'<div style="height:{{widthpx}}px;width:{{widthpx}}px;'+
						'border-radius:50%;position:relative;'+
						'box-shadow: 0px 0px 20px 0px black;'+//底圆加10px模糊的阴影，水平偏移，垂直偏移，阴影模糊，阴影突出，颜色
						'">'+  
				//用来遮挡旋转进入左边的半圆部分
				'<div style="height:{{widthpx}}px;width:{{widthpx}}px;position:absolute;clip:rect(0px,{{widthpx}}px,{{widthpx}}px,{{halfwidthpx}}px);">'+
					//可旋转的右半圆,这里的颜色将会是显示剩下percent值
					'<div id="rightcir" style="height: {{widthpx}}px;width:{{widthpx}}px;background-color:rgb(255,255,255);'+
								'border-radius:50%;'+
								'position:absolute;top:0;left:0;clip:rect(0px,{{widthpx}}px,{{widthpx}}px,{{halfwidthpx}}px);'+//遮挡左边
								'transform:rotate({{rightCirDeg}});'+//半圆裁剪，一定要position:absolute或fixed	
								'-o-transform:rotate({{rightCirDeg}});'+
								'-webkit-transform:rotate({{rightCirDeg}});'+
								'-moz-transform:rotate({{rightCirDeg}});'+
								'">'+
					'</div>'+
				'</div>'+
				//用来遮挡旋转进入右边的半圆部分（左半圆部分）
				'<div style="height:{{widthpx}}px;width:{{widthpx}}px;position:absolute;clip:rect(0px,{{halfwidthpx}}px,{{widthpx}}px,0px);">'+
					//可旋转的左半圆，颜色显示剩下percent
					'<div id="leftcir" style="height: {{widthpx}}px;width:{{widthpx}}px;background-color:rgb(255,255,255);'+
								'border-radius:50%;'+
								'position:absolute;top:0;left:0;clip:rect(0px,{{halfwidthpx}}px,{{widthpx}}px,0px);'+//遮挡左边
								'transform:rotate({{leftCirDeg}});'+//半圆裁剪，一定要position:absolute或fixed	
								'-o-transform:rotate({{leftCirDeg}});'+
								'-webkit-transform:rotate({{leftCirDeg}});'+
								'-moz-transform:rotate({{leftCirDeg}});'+
								'">'+
					'</div>'+
				'</div>'+
				//覆盖在正中心的圆，形成圆环，颜色是中间部分的颜色
				'<div style="height:80%;width:80%;background-color:rgb(255,255,255);'+
							'top:10%;left:10%;position:absolute;'+//计算与80%有关
							'border-radius:50%;'+
							'box-shadow: 0px 0px 10px 0px black inset;'+//内圆加10px模糊的内阴影
							'text-align:center;'+
							'">'+
					//字体有阴影
					//显示优,text-shadow:0px 0px 3px black;加了后有时刷新会看不清
					'<span style="font-size:{{fontsize1}}px;display:block;font-weight:bold;'+
								'text-shadow:0px 0px 3px black;'+
					'">优</span>'+
					//显示34℃
					//'<span style="font-size:{{fontsize2}}px;display:block;">{{stvalue}}{{unit}}</span>'+
					//显示图标+温度
					//'<span class="{{iconclass}}" style="font-size:{{iconsize}}px;display:block;">{{type}}</span>'+//图标需要先引入了字体css
				'</div>'+
			'</div>'+
			'</div>',
		link:function(scope,elem,attrs,ctrl){
			var fatherdiv=elem.find("div").eq(0);//主父div用来确定长宽的最小值
			var background=elem.find("div").eq(1);//背景底圆，也是占百分比的颜色条
			var leftcir=elem.find("div").eq(5);//左半圆
			var rightcir=elem.find("div").eq(3);//右半圆
			var statetext=elem.find("span").eq(0);//显示优差良的文字
			//=============初始化一些dom的值==============
			var init=function(){
				scope.widthpx=fatherdiv[0].clientWidth<fatherdiv[0].clientHeight ? fatherdiv[0].clientWidth : fatherdiv[0].clientHeight;
				scope.halfwidthpx=scope.widthpx/2;//自适应掩膜
				//自适应字体
				scope.fontsize1=(100/200*scope.widthpx+1).toFixed(0);//优字的大小，以200px下35px为标准
				//scope.iconsize=(30/200*scope.widthpx+1).toFixed(0);//
				//scope.fontsize2=(30/200*scope.widthpx+1).toFixed(0);
			}();//在最后加个括号可以让其马上运行,这样可以不加scope.apply()更新template
			//===========初始化结束===================

			//========监测value 的值改变css的一些显示====
			function updateDom(){
				var percent=(parseFloat(scope.stvalue)-parseFloat(scope.minStvalue))/(parseFloat(scope.maxStvalue)-parseFloat(scope.minStvalue));
				var percentInt=(percent*100).toFixed(0);//取整
				if (percentInt<=50 && percentInt>0) {//如果显示值小于50%，让右半园转动，左半圆不动
					scope.rightCirDeg=percentInt/100*360+'deg';
					scope.leftCirDeg='0deg';
					//根据百分比确定颜色,从百分0到百分百，绿色到红色渐变
					//小于50%时绿固定，红色慢慢增加，50%时红色和绿色都达到最大为橙色
					var percentcolred=(percentInt/100*255*2).toFixed(0);//百分比越多越红
					var percentcolgreed=255;
					background.css('background-color','rgb('+percentcolred+','+percentcolgreed+',0)');	
				}else if(percentInt<=100 && percentInt>50){
					scope.rightCirDeg='180deg';//右半圆旋转进不可见区，
					scope.leftCirDeg=(percentInt-50)/100*360+'deg';
					//大于百分50时红色最大维持不变，绿色慢慢减少，最后为纯红
					var percentcolred=255;//红色不变
					var percentcolgreed=((100-percentInt)/100*255*2).toFixed(0);
					background.css('background-color','rgb('+percentcolred+','+percentcolgreed+',0)');
				}else{
					alert("输入超出范围，见控制台");
					console.log("stvalue超出maxStvalue和minStvalue之间")
				};
				//优差良的文字显示，和对应颜色
				statetext.css('color',background.css('background-color'));
				if(percentInt<33){
					statetext.html("优");//小于百分33为优
				}else if(percentInt<66){
					statetext.html("良");
				}else{
					statetext.html("差"); 
				}
			}//();//在最后加个括号可以让其马上运行,这样可以不加scope.apply()更新template
			scope.$watch('stvalue',updateDom);
			//=====updatedom结束=======================
		}
	};
});
/*
楼主可以搜一下关于js拖动的代码，思路大体上就是： 
1.鼠标按下时（mousedown），将当前的滑块位置（css中top的值）和鼠标位置(在屏幕上的纵坐标位置pageY)记录下来，同时给滑块绑定鼠标移动事件函数； 
2.鼠标移动时（mousemove），随着移动事件的发生，绑定的移动函数会不断的执行（每次具体执行过程为：计算当前pageY和之前记录下的pageY的差值，然后将“差值+之前记录的top值”得到当前滑块应该在的top值，最后将top值写到滑块的样式中。），滑块因为top不断变化，上下位置也就会不断变化。 
3.鼠标松开时（mouseup），取鼠标消移动函数的绑定。
<sliderbar  wendu-md='wenduvalue'></sliderbar>
*/
model.directive('sliderbar',function(){
	//初始化控件的状态
	return{
		restrict:'EA',
		scope:{
			//按下鼠标时记录的top值
			//按下鼠标时的pagey值
			//slwidth:'@',
			wendu:'=wenduMd',
		},
		template:'<div style="background-color:rgb(230,230,230);height: 50%;width:100%;">'+
		'<h1 style="text-align:center;">目标温度:<span>{{wendu}}℃</span><h1>'+
				'</div>'+
				'<div style="background-color:rgb(230,230,23,0.4);height: 50%;width:100%;position:relative;">'+
					'<!-- 滑动条 --><div style="background-color:rgb(111,122,111);height: 50%;width:100%;top:25%;position:absolute;border-radius:10px;">'+
					'</div>'+
					'<!-- 滑动块 --><div class="block" style="left:0px;background-color:rgba(11,0,233,0.8);height: 80%;width:10%;top:10%;position:absolute;border-radius:10px;">'+
					'</div>'+
				'</div>',
		link:function(scope,elem,attrs,ctrl){
			var scoll=elem.find("div").eq(3);//eq(3)要根据上面的html变化而变化，指滑动块
			var scollline=elem.find("div").eq(2);//eq(2)要根据上面的html变化而变化，指滑动背景条
			var wendutext=elem.find("span").eq(0);//显示温度的文字
			var hasbind=false;//变量用了避免异步处理中还没unbind又bind了
			var curLeft;
			//初始化控件的状态
			
			//初始化控件的状态结束
			scoll.bind('mousedown',function(e){
				console.log("down");
				var downLeft=parseInt(scoll.css('left'));//转化为整型,滑块的当前高
				var downPageX=e.pageX;//按下鼠标时的pagey
				//按下时绑定移动事件
				if (!hasbind) {//只有当前没有绑定才绑定，避免了多次绑定
					console.log("bind");
					elem.bind('mousemove',a=function(e){//注意这里不用滑块的div而是用控件的整体div是避免鼠标移出滑块区却没响应
					 	hasbind=true;
					 	curLeft=e.pageX-downPageX+downLeft;
					 	console.log(e.pageX);
					 	//console.log(scollline[0].clientWidth);//用js获取宽度
					 	if (curLeft < scollline[0].clientWidth-scoll[0].clientWidth && curLeft > 0) {//限制滑块移动范围
					 		console.log(curLeft+"left");
					 		scoll.css('left',curLeft+'px');
					 		var leftpercent=curLeft/(scollline[0].clientWidth-scoll[0].clientWidth);//当前位位置占滑动条的百分比
					 		scope.wendu=parseFloat((leftpercent*20).toFixed(1))+15;//更新温度值,取小数点后2位
					 		var redcolor=(leftpercent*254).toFixed(0);//用来更新滑块的颜色
					 		var bluecolor=254-redcolor;
					 		scoll.css('background-color','rgba('+redcolor+',0,'+bluecolor+',0.8)');
					 		wendutext.css('color','rgba('+redcolor+',0,'+bluecolor+',0.8)');
					 		scope.$apply();//没了这个是不会在template里面更新的。
					 	};	
					});
				};
			});
			//松开鼠标时解除移动事件
			elem.bind('mouseup',function(e){
				hasbind=false;
				elem.unbind('mousemove',a);
				console.log("unbind");
			});
			//鼠标离开时也时解除移动事件

		}
	};
});
model.controller('page1mainCtrl',[
	'$scope',
	'$rootScope',
	function ($scope,$rootScope){
		$scope.wenduvalue="30";
		$scope.shiduvalue="58";
		$scope.guangzhaovalue="44";
		$scope.anqivalue="79";
	}
]);
