import React, { Component } from 'react';
import ExchangeChart from './ExchangeChart'; 
import {Navbar, NavbarBrand, Nav, NavItem, NavLink, TabPane, TabContent} from 'reactstrap';
import classnames from 'classnames';
import config from "./config.json";
import { MoonLoader } from 'halogenium';

class App extends Component {
  constructor(props) {
    super(props);

    this.toggle = this.toggle.bind(this);
    this.currencies = []
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

        this.setState(() => ({
          loading: false
        }));
      })
      .catch((err) => {
        alert(err)
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
            <MoonLoader color="#26A65B" size="64px" margin="4px" style={{position: "absolute", left: "50%", top: "50%"}}/>
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