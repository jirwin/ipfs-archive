/**
*
* ArchiveForm
*
*/

import React from 'react';
import PropTypes from 'prop-types';
// import styled from 'styled-components';
import { Row, Col, Button, Input } from 'reactstrap';
import FontAwesome from 'react-fontawesome';
import validUrl from 'valid-url';

class ArchiveForm extends React.PureComponent { // eslint-disable-line react/prefer-stateless-function
  constructor(props) {
    super(props);
    this.state = { url: '' };

    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  handleChange(event) {
    this.setState({ url: event.target.value });
  }

  handleSubmit(event) {
    event.preventDefault();
    if (validUrl.isWebUri(this.state.url)) {
      this.props.onSubmit(this.state.url);
    } else {
      this.props.onError(`${this.state.url} is not a valid url.`);
    }
  }

  button() {
    if (this.props.loading) {
      return (<Button className="form-control btn btn-success" disabled>
        <FontAwesome name="spinner" spin /> Archiving
      </Button>);
    }

    return <Button className="form-control btn btn-success" onClick={this.handleSubmit}>Archive</Button>;
  }

  render() {
    return (
      <div>
        <Row className="pb-3">
          <Col lg={10} md={10}>
            <Input onChange={this.handleChange} placeholder="Enter the URL you would like to archive..." />
          </Col>
          <Col>
            {this.button()}
          </Col>
        </Row>
      </div>
    );
  }
}

ArchiveForm.propTypes = {
  onSubmit: PropTypes.func.isRequired,
  onError: PropTypes.func.isRequired,
  loading: PropTypes.bool.isRequired,
};

export default ArchiveForm;
