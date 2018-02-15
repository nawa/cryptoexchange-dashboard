import React, { Component } from 'react';
import ExchangeChart from './ExchangeChart'; 
import {Table, Navbar, NavbarBrand, Nav, NavItem, NavLink, TabPane, TabContent} from 'reactstrap';
import classnames from 'classnames';
import config from "./config.json";
import { BeatLoader } from 'halogenium';
import moment from 'moment';

class App extends Component {
  constructor(props) {
    super(props);

    this.toggle = this.toggle.bind(this);
    this.currencies = []
    this.orders = []
    this.state = {
      activeTab: 'total',
      loading: true
    };
  }

  toggle(tab) {
    if (this.state.activeTab !== tab) {
      this.setState({
        activeTab: tab
      });
    }
  }

  componentWillMount() {
    return Promise.all([this.fetchActiveCurrencies(), this.fetchOrders()])
      .then(() => {
        this.setState(() => ({
          loading: false
        }));
      }).catch((err) => {
        alert(err)
      })
  }

  fetchActiveCurrencies() {
    return fetch(config.backend + "/balance/active", {
      method: 'GET'
    })
      .then((response) => response.json())
      .then((responseJson) => {
        this.currencies = Object.entries(responseJson)
          .sort((a, b) => {
            return b[1][0].btc - a[1][0].btc
          })
          .filter((a) => a[0] !== "total")
          .map((a) => a[0])
      })
  }

  fetchOrders() {
    return fetch(config.backend + "/order", {
      method: 'GET'
    })
      .then((response) => response.json())
      .then((responseJson) => {
        this.orders = responseJson
      })
  }

  render() {
    return (
      <div>
        <Navbar color="dark" dark>
          <NavbarBrand href="/">CryptoExchange Wallet Info</NavbarBrand>
          <Nav className="navbar-nav">
            <NavItem>
              <NavLink href="https://github.com/">Github</NavLink>
            </NavItem>
          </Nav>
        </Navbar>
        {
          (this.state.loading) 
            ?
            <BeatLoader color="#26A65B" size="32px" style={{position: "absolute", left: "50%", top: "50%"}}/>
            :
            <div>
              <Nav tabs>
                <NavItem style={{cursor: "pointer"}}>
                  <NavLink
                    className={classnames({ active: this.state.activeTab === 'total' })}
                    onClick={() => { this.toggle('total'); }}
                  >
                    BTC Total
                  </NavLink>
                </NavItem>
                {this.currencies.map((item, index) => (
                  <NavItem style={{cursor: "pointer"}}>
                    <NavLink
                      className={classnames({ active: this.state.activeTab === item })}
                      onClick={() => { this.toggle(item); }}
                    >
                      {item}
                    </NavLink>
                  </NavItem>
                ))}
              </Nav>
              <TabContent activeTab={this.state.activeTab}>
                <TabPane tabId="total">
                  <ExchangeChart currency="total" title="BTC Total"/>
                  <Table>
                    <thead>
                      <tr>
                      {/* //	|	Market	|	Time	|	BuyRate	| Can Sell=SellNowRate	|	Amount	|	Buy Price=BuyRate*Amount
                      //	|	Sell Price=SellNowPrice*Amount	|	Profit = BuyPrice+BuyPrice*0.0025 - (Sell Price - Sell Price*0.0025)
                      //	|	Profit BTC = Profit * BTCRate |	Profit USDT = Profit * USDTRate | Profit % = (Profit / Amount)*100 */}
                        <th>Market</th>
                        <th>Date</th>
                        <th>Amount</th>
                        <th>Buy Rate</th>
                        <th>Can Sell Rate</th>
                        <th>Buy Price</th>
                        <th>Sell Price</th>
                        <th>Profit %</th>
                        <th>Profit</th>
                        <th>Profit USDT</th>
                      </tr>
                    </thead>
                    <tbody>
                      {this.orders.map((order, index) => {
                        const buyPrice = order.amount * order.buy_rate
                        const sellPrice = order.amount * order.sellnow_rate
                        const profit = (sellPrice - sellPrice * 0.0025) - (buyPrice + buyPrice * 0.0025)
                        return (
                          <tr>
                            <td><a href={order.market_link} target="_blank">{order.market}</a></td>
                            <td>{moment.unix(order.time).format('YY-MM-DD HH:mm')}</td>
                            <td>{order.amount}</td>
                            <td>{order.buy_rate.toFixed(8)}</td>
                            <td>{order.sellnow_rate.toFixed(8)}</td>
                            <td>{(order.amount * order.buy_rate).toFixed(8)}</td>
                            <td>{sellPrice.toFixed(8)}</td>
                            <th style={profit > 0 ? {color: "green"} : {color: "red"}}>{(profit / buyPrice * 100).toFixed(2)}%</th>
                            <th style={profit > 0 ? {color: "green"} : {color: "red"}}>{profit.toFixed(8)}</th>
                            <th style={profit > 0 ? {color: "green"} : {color: "red"}}>{(profit * order.usdt_rate).toFixed(2)}</th>
                          </tr>
                        )
                      })}
                    </tbody>
                  </Table>
                </TabPane>
                {this.currencies.map((item, index) => (
                  <TabPane tabId={item}>
                    <ExchangeChart currency={item} title={item}/>
                  </TabPane>
                ))}
              </TabContent>
            </div>
        }
      </div>
    );
  }
}

export default App;