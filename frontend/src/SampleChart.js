import React from 'react';
import Label, {LineChart, Line, CartesianGrid, XAxis, YAxis, Tooltip, ReferenceArea} from 'recharts';

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
  // top : dataMax => {
  //   return (dataMax + Math.abs(dataMax)*0.2)
  // },
  // bottom : dataMin => {
  //   (dataMin - Math.abs(dataMin)*0.2)
  // },
  // top2 : dataMax => {
  //   (dataMax + Math.abs(dataMax)*0.2)
  // },
  // bottom2 : dataMin => {
  //   (dataMin - Math.abs(dataMin)*0.2)
  // },
  top : 'auto',
  bottom : 'auto',
  top2 : 'auto',
  bottom2 : 'auto',
  animation : true
};

class ExchangeChart extends React.Component {

	constructor(props) {
    super(props);
    this.state = initialState;
  }
  
  zoom(){  
  	let { refAreaLeft, refAreaRight, data } = this.state;

		if ( refAreaLeft === refAreaRight || refAreaRight === '' ) {
    	this.setState( () => ({
      	refAreaLeft : '',
        refAreaRight : ''
      }) );
    	return;
    }

		// xAxis domain
	  if ( refAreaLeft > refAreaRight ) 
    		[ refAreaLeft, refAreaRight ] = [ refAreaRight, refAreaLeft ];

    // yAxis domain
    const from = this.state.data.findIndex((v) => { 
      return v.name == refAreaLeft 
    })
    const to = this.state.data.findIndex((v) => { 
      return v.name == refAreaRight
    })
    const [ bottom, top ] = getAxisYDomain( this.state.data, from, to, 'cost', 0 );
    const [ bottom2, top2 ] = getAxisYDomain( this.state.data, from, to, 'impression', 0);

    this.setState( () => ({
      refAreaLeft : '',
      refAreaRight : '',
    	data : data.slice(),
      left : refAreaLeft,
      right : refAreaRight,
      bottom, top, bottom2, top2
    } ) );
  }

	zoomOut() {
  	const { data } = this.state;
  	this.setState( () => ({
      data : data.slice(),
      refAreaLeft : '',
      refAreaRight : '',
      left : 'dataMin',
      right : 'dataMax',
      top : 'auto',
      bottom : 'auto',
      top2 : 'auto',
      bottom2: 'auto'
    }) );
  }

  componentDidMount() {
    // this.updateData(data)
  }

  updateData(data) {
    const data1 = data;
    this.setState( () => ({
      data : data1.slice()
    }));
  }
  
  render() {
    const { data, barIndex, left, right, refAreaLeft, refAreaRight, top, bottom, top2, bottom2 } = this.state;

    return (
      <div className="highlight-bar-charts">
        <a
          href="javascript: void(0);"
          className="btn update"
          onClick={this.zoomOut.bind( this )}
        >
          Zoom Out
        </a>


        <p>Highlight / Zoom - able Line Chart</p>
          <LineChart
            width={800}
            height={400}
            data={data}
            onMouseDown = { 
              (e) => this.setState({refAreaLeft:e.activeLabel}) 
            }
            onMouseMove = { 
              (e) => {
                this.state.refAreaLeft && this.setState({refAreaRight:e.activeLabel})
              }
            }
            onMouseUp = { this.zoom.bind( this ) }
          >
            <CartesianGrid strokeDasharray="3 3"/>
            <XAxis 
              allowDataOverflow={true}
              dataKey="name"
              domain={[left, right]}
              type="number"
            />
            <YAxis 
              allowDataOverflow={true}
              domain={[bottom, top]}
              type="number"
              yAxisId="1"
             />
            <YAxis 
              orientation="right"
              allowDataOverflow={true}
              domain={[bottom2, top2]}
              type="number"
              yAxisId="2"
             /> 
            <Tooltip/>
            <Line yAxisId="1" type='natural' dataKey='cost' stroke='#8884d8' animationDuration={300}/>
            <Line yAxisId="2" type='natural' dataKey='impression' stroke='#82ca9d' animationDuration={300}/>
            
            {
            	(refAreaLeft && refAreaRight) ? (
              <ReferenceArea yAxisId="1" x1={refAreaLeft} x2={refAreaRight}  strokeOpacity={0.3} /> ) : null
            
            }
            
          </LineChart> 

      </div>
    );
  }
}

export default ExchangeChart;