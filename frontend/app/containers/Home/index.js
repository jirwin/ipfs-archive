/**
 *
 * Home
 *
 */

import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { createStructuredSelector } from 'reselect';
import { compose } from 'redux';

import { Navbar, NavbarBrand, Nav, NavItem, NavLink, Jumbotron, Container, Row, Col } from 'reactstrap';

import injectSaga from 'utils/injectSaga';
import injectReducer from 'utils/injectReducer';
import makeSelectHome from './selectors';
import reducer from './reducer';
import saga from './saga';
import {
  archiveAction,
  archiveError,
} from './actions';

import ArchiveForm from '../../components/ArchiveForm';

class Error extends React.PureComponent {
  render() {
    if (this.props.msg !== '') {
      return (
        <div className="alert alert-danger">{this.props.msg}</div>
      );
    }

    return null;
  }
}

Error.propTypes = {
  msg: PropTypes.string.isRequired,
};

class DisplayResults extends React.PureComponent {
  render() {
    return (
      <div className="alert alert-success">{this.props.url} is archived at <a className="alert-link" href={this.props.archivedUrl}>{this.props.archivedUrl}</a></div>
    );
  }
}

DisplayResults.propTypes = {
  url: PropTypes.string.isRequired,
  archivedUrl: PropTypes.string.isRequired,
};

export class Home extends React.Component { // eslint-disable-line react/prefer-stateless-function
  getContent() {
    if (this.props.home.error) {
      return (
        <Error msg={this.props.home.error} />
      );
    }

    if (this.props.home.archivedUrl && this.props.home.url) {
      return (
        <DisplayResults url={this.props.home.url} archivedUrl={this.props.home.archivedUrl} />
      );
    }

    return null;
  }

  render() {
    return (
      <Container>
        <Navbar color="dark" dark expand>
          <NavbarBrand href="/">ipfs archive network</NavbarBrand>
          <Nav className="ml-auto" navbar>
            <NavItem>
              <NavLink href="https://github.com/jirwin/ipfs-archive">github</NavLink>
            </NavItem>
          </Nav>
        </Navbar>
        <Jumbotron>
          <Row>
            <Col>
              <h1 className="display-3">save your stuff</h1>
              <p>
                    use the ipfs archive network to create a snapshot of any url.
                  </p>
            </Col>
          </Row>
        </Jumbotron>
        <ArchiveForm loading={this.props.home.loading} onSubmit={this.props.onArchiveSubmit} onError={this.props.onArchiveError} />
        <Row>
          <Col>
            {this.getContent()}
          </Col>
        </Row>
      </Container>
    );
  }
}

Home.propTypes = {
  dispatch: PropTypes.func.isRequired,
};

const mapStateToProps = createStructuredSelector({
  home: makeSelectHome(),
});

function mapDispatchToProps(dispatch) {
  return {
    onArchiveSubmit: (url) => {
      dispatch(archiveAction(url));
    },
    onArchiveError: (err) => {
      dispatch(archiveError(err));
    },
    dispatch,
  };
}

const withConnect = connect(mapStateToProps, mapDispatchToProps);

const withReducer = injectReducer({ key: 'home', reducer });
const withSaga = injectSaga({ key: 'home', saga });

export default compose(
  withReducer,
  withSaga,
  withConnect,
)(Home);
