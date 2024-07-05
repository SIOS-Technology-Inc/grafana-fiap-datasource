import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { userEvent } from '@testing-library/user-event';

import { ConfigEditor } from './ConfigEditor';

describe('ConfigEditor', () => {
  const testOptions = {
    "id": 2,
    "uid": "a25326a4-956f-4006-8b21-1e4c3544772d",
    "orgId": 1,
    "name": "Fiap",
    "type": "sios-fiap-datasource",
    "typeName": "A",
    "typeLogoUrl": "",
    "access": "proxy",
    "url": "",
    "user": "",
    "database": "",
    "basicAuth": false,
    "basicAuthUser": "",
    "withCredentials": false,
    "isDefault": true,
    "jsonData": {
      "url": ""
    },
    "secureJsonFields": {},
    "version": 25,
    "readOnly": false,
  }

  const onOptionsChange = jest.fn();

  describe('url validation test', () => {
    describe('when url is empty in the initial state', () => {
      it('should show error message', async () => {
        render(<ConfigEditor onOptionsChange={onOptionsChange} options={testOptions} />);

        await waitFor(() => {
          expect(screen.queryByText('This field is required')).toBeInTheDocument();
        });
      });
    });
    describe('when url is empty after input', () => {
      it('should show error message', async () => {
        render(<ConfigEditor onOptionsChange={onOptionsChange} options={testOptions} />);

        const input = screen.getByRole('textbox', { name: /url/i });
        await userEvent.type(input, 'http://test.server.com:8080');
        await userEvent.clear(input);

        await waitFor(() => {
          expect(screen.queryByText('This field is required')).toBeInTheDocument();
        });
      });
    });
    describe('when url is invalid', () => {
      const InputURLs = ['http://', 'https://', 'htt://a', ' http://a',' https://a']
      it.each(InputURLs)('should show error message (input: %s)', async (inputURL) => {
        render(<ConfigEditor onOptionsChange={onOptionsChange} options={testOptions} />);

        const input = screen.getByRole('textbox', { name: /url/i });

        await userEvent.type(input, inputURL);
        
        await waitFor(() => {
          expect(screen.queryByText('Invalid URL format.')).toBeInTheDocument();
        });
      });
    });
    describe('when url is valid', () => {
      const InputURLs = ['http://a', 'https://a']
      it.each(InputURLs)('should show error message (input: %s)', async (inputURL) => {
        render(<ConfigEditor onOptionsChange={onOptionsChange} options={testOptions} />);

        const input = screen.getByRole('textbox', { name: /url/i });

        await userEvent.type(input, inputURL);
        
        await waitFor(() => {
          expect(screen.queryByText('Invalid URL format.')).not.toBeInTheDocument();
        });
      });
    });
  });
});
