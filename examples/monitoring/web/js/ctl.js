function Ctl() {
  this.cache_num = 16;
  this.metric_names = ["memory_used", "cpu_usage", "eth0_rx", "eth0_tx", "eth1_rx", "eth1_tx"];
  this.agents = {}
}

Ctl.prototype.init_agent = function(agentName) {
  var metrics = {}
  for(var i = 0; i < this.metric_names.length; i++) {
    metrics[this.metric_names[i]] = {
      next: 0,
      data: [] // circular buffer
    };
    for(var j=0; j < this.cache_num; j++) {
      metrics[this.metric_names[i]].data.push({
        time: new Date().getTime(),
        value: 0
      })
    }
  }
  this.agents[agentName] = metrics;
  this.new_agent_graph(agentName);
}

Ctl.prototype.poll = function(clientID){
  var self = this;
  d3.json("/data?clientID=" + clientID, function(data) {
    if(!self.agents[data.AgentName]) {
      self.init_agent(data.AgentName);
    }
    self.agents[data.AgentName][data.Name].data[self.agents[data.AgentName][data.Name].next] = {
      time: new Date().getTime(),
      value: data.Data
    };
    self.agents[data.AgentName][data.Name].next = index_plus(self.agents[data.AgentName][data.Name].next, self.cache_num, 1);
    self.poll(data.ClientID);
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

Ctl.prototype.get_metric = function(agent, metric_name) {
  var self = this;
  return self.context.metric(function(start, stop, step, callback) {
    var num = (stop - start) / step;
    if(num > self.cache_num) {
      callback(new Error("Too many values requested"));
      return;
    }
    var last = index_minus(agent[metric_name].next, self.cache_num, 1);
    for (var i = 0; agent[metric_name].data[last].time > stop; i++) {
      if(i == self.cache_num) {
        callback(new Error("Data for the time period is not available"));
        return;
      }
      last = index_minus(last, self.cache_num, 1);
    }
    var ret = [];
    for (var i = index_minus(last, self.cache_num, num); ret.length < num; i = index_plus(i, self.cache_num, 1)) {
      ret.push(agent[metric_name].data[i].value);
    }
    callback(null, ret);
  }, metric_name);
}

Ctl.prototype.new_agent_graph = function(agentName) {
  var self = this;
  var cubism_metrics = [];
  for(var i = 0; i < self.metric_names.length; i++) {
    cubism_metrics.push(self.get_metric(self.agents[agentName], self.metric_names[i]));
  }

  d3.select("#agents").call(function(agentDiv) {

    var div = agentDiv.append("div").attr("id", agentName);

    div.append("div").attr("class", "agent_title").text(agentName);

    div.append("div")
    .attr("class", "axis")
    .call(self.context.axis().orient("top"));

    div.selectAll(".horizon")
    .data(cubism_metrics)
    .enter().append("div")
    .attr("class", "horizon")
    .call(self.context.horizon().extent(null));

    div.append("div")
    .attr("class", "rule")
    .call(self.context.rule());

  });

  self.context.on("focus", function(i) {
    d3.selectAll(".value").style("right", i == null ? null : self.context.size() - i + "px");
  });

}

Ctl.prototype.config_cubism = function(config) {
  this.context = cubism.context()
  .serverDelay(0)
  .clientDelay(0)
  .step(config.interval)
  .size(960);
}
