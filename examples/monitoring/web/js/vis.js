function newCtx() {
  var ctx = {}
  ctx.cache_num = 16;
  ctx.metric_names = ["memory_used", "cpu_usage", "eth0_rx", "eth0_tx", "eth1_rx", "eth1_tx"];
  ctx.agents = {}
  return ctx;
}

function init_agent(ctx, agentName) {
  var metrics = {}
  for(var i = 0; i < ctx.metric_names.length; i++) {
    metrics[ctx.metric_names[i]] = {
      next: 0,
      data: [] // circular buffer
    };
    for(var j=0; j < ctx.cache_num; j++) {
      metrics[ctx.metric_names[i]].data.push({
        time: new Date().getTime(),
        value: 0
      })
    }
  }
  ctx.agents[agentName] = metrics;
  new_agent_graph(ctx, agentName)
}

function poll(ctx, clientID){
  d3.json("/data?clientID=" + clientID, function(data) {
    if(!ctx.agents[data.AgentName]) {
      init_agent(ctx, data.AgentName);
    }
    ctx.agents[data.AgentName][data.Name].data[ctx.agents[data.AgentName][data.Name].next] = {
      time: new Date().getTime(),
      value: data.Data
    };
    ctx.agents[data.AgentName][data.Name].next = index_plus(ctx.agents[data.AgentName][data.Name].next, ctx.cache_num, 1);
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

function get_metric(ctx, agent, metric_name) {
  return ctx.context.metric(function(start, stop, step, callback) {
    var num = (stop - start) / step;
    if(num > ctx.cache_num) {
      callback(new Error("Too many values requested"));
      return;
    }
    var last = index_minus(agent[metric_name].next, ctx.cache_num, 1);
    for (var i = 0; agent[metric_name].data[last].time > stop; i++) {
      if(i == ctx.cache_num) {
        callback(new Error("Data for the time period is not available"));
        return;
      }
      last = index_minus(last, ctx.cache_num, 1);
    }
    var ret = [];
    for (var i = index_minus(last, ctx.cache_num, num); ret.length < num; i = index_plus(i, ctx.cache_num, 1)) {
      ret.push(agent[metric_name].data[i].value);
    }
    callback(null, ret);
  }, metric_name);
}

function new_agent_graph(ctx, agentName) {
  var cubism_metrics = [];
  for(var i = 0; i < ctx.metric_names.length; i++) {
    cubism_metrics.push(get_metric(ctx, ctx.agents[agentName], ctx.metric_names[i]));
  }

  d3.select("#agents").call(function(agentDiv) {
    
    var div = agentDiv.append("div").attr("id", agentName);

    div.append("div").attr("class", "agent_title").text(agentName);

    div.append("div")
    .attr("class", "axis")
    .call(ctx.context.axis().orient("top"));

    div.selectAll(".horizon")
    .data(cubism_metrics)
    .enter().append("div")
    .attr("class", "horizon")
    .call(ctx.context.horizon().extent(null));

    div.append("div")
    .attr("class", "rule")
    .call(ctx.context.rule());

  });

  ctx.context.on("focus", function(i) {
    d3.selectAll(".value").style("right", i == null ? null : ctx.context.size() - i + "px");
  });

}

function config_cubism(ctx, config) {
  ctx.context = cubism.context()
  .serverDelay(0)
  .clientDelay(0)
  .step(config.interval)
  .size(960);
}

function run() {
  var ctx = newCtx();
  d3.json("/config", function(config) {
    config_cubism(ctx, config);
    poll(ctx, "");
  });
}

run();
