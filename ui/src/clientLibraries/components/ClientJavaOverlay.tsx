// Libraries
import React, {FunctionComponent} from 'react'
import {connect} from 'react-redux'

// Components
import ClientLibraryOverlay from 'src/clientLibraries/components/ClientLibraryOverlay'
import TemplatedCodeSnippet from 'src/shared/components/TemplatedCodeSnippet'

// Constants
import {clientJavaLibrary} from 'src/clientLibraries/constants'

// Types
import {AppState} from 'src/types'

interface StateProps {
  org: string
}

type Props = StateProps

const ClientJavaOverlay: FunctionComponent<Props> = props => {
  const {
    name,
    url,
    buildWithMavenCodeSnippet,
    buildWithGradleCodeSnippet,
    initializeClientCodeSnippet,
    executeQueryCodeSnippet,
    writingDataLineProtocolCodeSnippet,
    writingDataPointCodeSnippet,
    writingDataPojoCodeSnippet,
    pojoClassCodeSnippet,
  } = clientJavaLibrary
  const {org} = props
  const server = window.location.origin

  return (
    <ClientLibraryOverlay title={`${name} Client Library`}>
      <p>
        For more detailed and up to date information check out the{' '}
        <a href={url} target="_blank">
          GitHub Repository
        </a>
      </p>
      <h5>Add Dependency</h5>
      <p>Build with Maven</p>
      <TemplatedCodeSnippet template={buildWithMavenCodeSnippet} label="Code" />
      <p>Build with Gradle</p>
      <TemplatedCodeSnippet
        template={buildWithGradleCodeSnippet}
        label="Code"
      />
      <h5>Initialize the Client</h5>
      <TemplatedCodeSnippet
        template={initializeClientCodeSnippet}
        label="Java Code"
        defaults={{
          server: 'serverUrl',
          token: 'token',
        }}
        values={{
          server,
        }}
      />
      <h5>Write Data</h5>
      <p>Option 1: Use InfluxDB Line Protocol to write data</p>
      <TemplatedCodeSnippet
        template={writingDataLineProtocolCodeSnippet}
        label="Java Code"
        defaults={{
          bucket: 'bucketID',
          org: 'orgID',
        }}
        values={{
          org,
        }}
      />
      <p>Option 2: Use a Data Point to write data</p>
      <TemplatedCodeSnippet
        template={writingDataPointCodeSnippet}
        label="Java Code"
        defaults={{
          bucket: 'bucketID',
          org: 'orgID',
        }}
        values={{
          org,
        }}
      />
      <p>Option 3: Use POJO and corresponding class to write data</p>
      <TemplatedCodeSnippet
        template={writingDataPojoCodeSnippet}
        label="Java Code"
        defaults={{
          bucket: 'bucketID',
          org: 'orgID',
        }}
        values={{
          org,
        }}
      />
      <TemplatedCodeSnippet template={pojoClassCodeSnippet} label="Java Code" />
      <h5>Execute a Flux query</h5>
      <TemplatedCodeSnippet
        template={executeQueryCodeSnippet}
        label="Java Code"
        defaults={{
          bucket: 'my_bucket',
          org: 'myorgid',
        }}
        values={{
          org,
        }}
      />
    </ClientLibraryOverlay>
  )
}

const mstp = (state: AppState): StateProps => {
  const org = state.orgs.org.id

  return {
    org,
  }
}

export {ClientJavaOverlay}
export default connect<StateProps, {}, Props>(
  mstp,
  null
)(ClientJavaOverlay)
