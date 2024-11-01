// force timezone to UTC to allow tests to work regardless of local timezone
// generally used by snapshots, but can affect specific tests
process.env.TZ = 'UTC';

module.exports = {
  // Use jsdom as test environment for use react-testing-library
  testEnvironment: 'jsdom',
  // Jest configuration provided by Grafana scaffolding
  ...require('./.config/jest.config'),
};
