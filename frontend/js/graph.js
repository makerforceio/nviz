var lossChart;

var graphInit = function () {
    var ctx = document.querySelector(".lossChart").getContext('2d');
    var data ={
        labels: [],
        datasets: [
            {
                label: "Loss",
                fill: "false",
                lineTension: 0,
                backgroundColor: "rgba(75,192,192,0.4)",
                borderColor: "rgba(75,192,192,1)",
                borderCapStyle: 'butt',
                borderDash: [0],
                borderDashOffset: 0.0,
                borderJoinStyle: 'miter',
                pointBorderColor: "rgba(75,192,192,1)",
                pointBackgroundColor: "#fff",
                pointBorderWidth: 1,
                pointHoverRadius: 5,
                pointHoverBackgroundColor: "rgba(75,192,192,1)",
                pointHoverBorderColor: "rgba(220,220,220,1)",
                pointHoverBorderWidth: 2,
                pointRadius: 1,
                pointHitRadius: 10,
                data: [],
                spanGaps: false,
            }
        ]
    };

    lossChart = new Chart(ctx, {
        type: 'line',
        data: data,
        options: {
			animationSteps: 15,
			scales: {
				xAxes: [{display: false}]
				yAxes: [{display: false}]
			},
			tooltips: {enabled: false}, legend: {display: false}
		}
    });
};

var graphAddDatapoint = function (data){
    lossChart.data.datasets[0].data.push(data);
    lossChart.update();
};

