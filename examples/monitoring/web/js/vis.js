function newCtx() {
  var ctx = {}
  ctx.cache_num = 16;
  ctx.metric_names = ["memory_used", "cpu_usage", "eth0_rx", "eth0_tx", "eth1_rx", "eth1_tx"];
  ctx.metrics = {};
  for(var i = 0; i < ctx.metric_names.length; i++) {
    ctx.metrics[ctx.metric_names[i]] = {
      next: 0,
      data: [] // circular buffer
    };
    for(var j=0; j < ctx.cache_num; j++) {
      ctx.metrics[ctx.metric_names[i]].data.push({
        time: new Date().getTime(),
        value: 0
      })
    }
  }
  return ctx;
}

function poll(ctx, clientID){
  d3.json("/data?clientID=" + clientID, function(data) {
    ctx.metrics[data.Name].data[ctx.metrics[data.Name].next] = {
      time: new Date().getTime(),
      value: data.Data
    };
    ctx.metrics[data.Name].next = index_plus(ctx.metrics[data.Name].next, ctx.cache_num, 1);
    poll(ctx, data.ClientID);
  });
}

function index_minus(original, count, offset) {
  if(offset >= count) {
    return -1
  }
  return original - offset >= 0 ? (original - offset) : (original - offset + count);
}

function index_plus(original, count, offset) {
  if(offset >= count) {
    return -1
  }
  return original + offset < count ? (original + offset) : (original + offset - count);
}

function get_metric(ctx, context, metric_name) {
  return context.metric(function(start, stop, step, callback) {
    var num = (stop - start) / step;
    if(num > ctx.cache_num) {
      callback(new Error("Too many values requested"));
      return;
    }
    var last = index_minus(ctx.metrics[metric_name].next, ctx.cache_num, 1);
    for (var i = 0; ctx.metrics[metric_name].data[last].time > stop; i++) {
      if(i == ctx.cache_num) {
        callback(new Error("Data for the time period is not available"));
        return;
      }
      last = index_minus(last, ctx.cache_num, 1);
    }
    var ret = [];
    for (var i = index_minus(last, ctx.cache_num, num); ret.length < num; i = index_plus(i, ctx.cache_num, 1)) {
      ret.push(ctx.metrics[metric_name].data[i].value);
    }
    callback(null, ret);
  }, metric_name);
}

function configCubism(ctx, config) {
  var context = cubism.context()
  .serverDelay(0)
  .clientDelay(0)
  .step(config.interval)
  .size(960);

  var cubism_metrics = [];
  for(var i = 0; i < ctx.metric_names.length; i++) {
    cubism_metrics.push(get_metric(ctx, context, ctx.metric_names[i]));
  }

  d3.select("#agent").call(function(div) {

    div.append("div")
    .attr("class", "axis")
    .call(context.axis().orient("top"));

    div.selectAll(".horizon")
    .data(cubism_metrics)
    .enter().append("div")
    .attr("class", "horizon")
    .call(context.horizon().extent(null));

    div.append("div")
    .attr("class", "rule")
    .call(context.rule());

  });

  context.on("focus", function(i) {
    d3.selectAll(".value").style("right", i == null ? null : context.size() - i + "px");
  });
}

function run() {
  var ctx = newCtx();
  poll(ctx, "");
  d3.json("/config", function(config) {
    configCubism(ctx, config);
  });
}

run();
