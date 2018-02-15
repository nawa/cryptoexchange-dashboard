import React from 'react';
import moment from 'moment';
import { LineChart, Line, CartesianGrid, XAxis, YAxis, Tooltip, ReferenceArea } from 'recharts';
import { Container, Row, Col, ButtonGroup, Button } from "reactstrap";
import config from "./config.json";
import { ScaleLoader } from 'halogenium';

const getAxisYDomain = (data, from, to, ref, offset) => {
  const refData = data.slice(from, to+1);
  let [ bottom, top ] = [ refData[0][ref], refData[0][ref] ];
  refData.forEach( d => {
  	if ( d[ref] > top ) top = d[ref];
    if ( d[ref] < bottom ) bottom = d[ref];
  });
  
  return [ (bottom|0) - offset, (top|0) + offset ]
};

const initialState = {
  data : [],
  left : 'dataMin',
  right : 'dataMax',
  refAreaLeft : '',
  refAreaRight : '',
  top : 'auto',
  bottom : 'auto',
  top2 : 'auto',
  bottom2 : 'auto',
  animation : false,
  zoom: false,
  period: "2h",
  loading: false
};

class ExchangeChart extends React.Component {

	constructor(props) {
    super(props);
    this.state = initialState;
    this.title = this.props.title
    this.currency = this.props.currency
  }

  componentDidMount() {
    this.load2h()
  }
  
  load2h() {
    return this.load('/balance/period/hourly/2?currency=' + this.currency)
      .then(() => {
        this.setState(() => ({
          period: "2h"
        }));
      })
  }

  load1d() {
    return this.load('/balance/period/hourly/24?currency=' + this.currency)
      .then(() => {
        this.setState(() => ({
          period: "1d"
        }));
      })
  }

  load1w() {
    return this.load('/balance/period/weekly?currency=' + this.currency)
      .then(() => {
        this.setState(() => ({
          period: "1w"
        }));
      })
  }

  load1m() {
    return this.load('/balance/period/monthly?currency=' + this.currency)
      .then(() => {
        this.setState(() => ({
          period: "1m"
        }));
      })
  }

  loadAll() {
    return this.load('/balance/period/all?currency=' + this.currency)
      .then(() => {
        this.setState(() => ({
          period: "all"
        }));
      })
  }

  load(endpoint) {
    this.setState(() => ({
      loading: true
    }));
    return fetch(config.backend + endpoint, {
      method: 'GET'
    })
      .then((response) => response.json())
      .then((responseJson) => {
        if (this.state.zoom) {
          this.zoomOut()
        }
        this.setState(() => ({
          data: responseJson[this.currency].slice(),
          loading : false
        }))
      })
      .catch((err) => {
        alert(err)
      })
  }
  
  zoom() {
    let { refAreaLeft, refAreaRight, data } = this.state;

    if (refAreaLeft === refAreaRight || refAreaRight === '') {
      this.setState(() => ({
        refAreaLeft: '',
        refAreaRight: ''
      }));
      return;
    }

    // xAxis domain
    if (refAreaLeft > refAreaRight)
      [refAreaLeft, refAreaRight] = [refAreaRight, refAreaLeft];

    // yAxis domain
    const from = this.state.data.findIndex((v) => {
      return v.time === refAreaLeft
    })
    const to = this.state.data.findIndex((v) => {
      return v.time === refAreaRight
    })

    const [bottom, top] = getAxisYDomain(this.state.data, Math.min(from, to), Math.max(from, to), 'usdt');
    const [bottom2, top2] = getAxisYDomain(this.state.data, Math.min(from, to), Math.max(from, to), 'btc');

    this.setState(() => ({
      refAreaLeft: '',
      refAreaRight: '',
      data: data.slice(),
      left: refAreaLeft,
      right: refAreaRight,
      zoom: true,
      bottom, top, bottom2, top2
    }));
  }

  zoomOut() {
    const { data } = this.state;
    this.setState(() => ({
      data: data.slice(),
      refAreaLeft: '',
      refAreaRight: '',
      left: 'dataMin',
      right: 'dataMax',
      top: 'auto',
      bottom: 'auto',
      top2: 'auto',
      bottom2: 'auto',
      zoom: false
    }));
  }

  timeFormat(time) {
    return moment.unix(time).format('YY-MM-DD HH:mm')
  }

  tooltipLabelFormatter(time) {
    return moment.unix(time).format('YYYY-MM-DD HH:mm:s')
  }

  render() {
    const { data, left, right, refAreaLeft, refAreaRight, top, bottom, top2, bottom2, zoom, period } = this.state;

    return (
      <Container>
        <br/>
        <Row>
          <Col>BTC - {data[0] ? data[0].btc : ""}</Col>
        </Row>
        <Row>
          <Col>USDT - {data[0] ? data[0].usdt : ""}</Col>
        </Row>
        <br/>
        <Row>
          <Col>
            <ButtonGroup>
              <Button active={period === "2h"} onClick={this.load2h.bind(this)}>2h</Button>
              <Button active={period === "1d"} onClick={this.load1d.bind(this)}>1d</Button>
              <Button active={period === "1w"} onClick={this.load1w.bind(this)}>1w</Button>
              <Button active={period === "1m"} onClick={this.load1m.bind(this)}>1m</Button>
              <Button active={period === "all"} onClick={this.loadAll.bind(this)}>All</Button>
            </ButtonGroup>
          </Col>
          <Col>{this.title}</Col>
          <Col>
            <Button disabled={!zoom} onClick={this.zoomOut.bind(this)}>Zoom Out</Button>
          </Col>
        </Row>
        <Row>
          <Col>
          {
            (this.state.loading) 
              ?
              <div style={{width: "1000px", height: "500px"}}>
                <ScaleLoader color="#26A65B" size="64px" style={{"padding-left": "470px", "padding-top": "220px"}}/>
              </div>
              :
              <LineChart
                width={1000}
                height={500}
                data={data}
                onMouseDown = {(e) => {
                  if (e) {
                    this.setState({refAreaLeft:e.activeLabel})
                  }
                }}
                onMouseMove = {(e) => this.state.refAreaLeft && this.setState({refAreaRight:e.activeLabel})}
                onMouseUp = { this.zoom.bind( this ) }
              >
                <CartesianGrid strokeDasharray="3 3"/>
                <XAxis 
                  allowDataOverflow={true}
                  dataKey="time"
                  domain={[left, right]}
                  tickFormatter={this.timeFormat}
                  type="number"
                  style={{"font-size": "x-small"}}
                  tickCount="10"
                />
                <YAxis 
                  allowDataOverflow={true}
                  domain={[bottom, top]}
                  type="number"
                  yAxisId="1"
                  tickCount="10"
                />
                <YAxis
                  orientation="right"
                  allowDataOverflow={true}
                  domain={[bottom2, top2]}
                  type="number"
                  yAxisId="2"
                  tickCount="10"
                /> 
                <Tooltip labelFormatter={this.tooltipLabelFormatter}/>
                <Line yAxisId="1" type='linear' dot={false} dataKey='usdt' stroke='#009e73' isAnimationActive={false} strokeWidth="2" />
                <Line yAxisId="2" type='linear' dot={false} dataKey='btc' stroke='#ff9300' isAnimationActive={false} strokeWidth="2" />
                
                {
                  (refAreaLeft && refAreaRight) ? (
                  <ReferenceArea yAxisId="1" x1={refAreaLeft} x2={refAreaRight}  strokeOpacity={0.3} /> ) : null
                }
                
              </LineChart>
          }
          </Col>
        </Row>
      </Container>
    );
  }
}

export default ExchangeChart;