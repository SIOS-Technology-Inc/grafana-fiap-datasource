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
        render(
          <ConfigEditor
            onOptionsChange={onOptionsChange}
            options={testOptions}
          />
        );

        waitFor(() => {
          expect(screen.getByText('This field is required')).toBeInTheDocument();
        });
      });
    });
    describe('when url is empty after input', () => {
      it('should show error message', async () => {
        render(
          <ConfigEditor
            onOptionsChange={onOptionsChange}
            options={testOptions}
          />
        );

        const input = screen.getByRole('textbox', { name: /url/i });
        userEvent.type(input, 'http://test.server.com:8080');
        userEvent.clear(input);

        waitFor(() => {
          expect(screen.getByText('This field is required')).toBeInTheDocument();
        });
      });
    });
    describe('when url is invalid', () => {
      it('should show error message', async () => {
        render(
          <ConfigEditor
            onOptionsChange={onOptionsChange}
            options={testOptions}
          />
        );

        render(<ConfigEditor onOptionsChange={onOptionsChange} options={testOptions} />);

        const input = screen.getByRole('textbox', { name: /url/i });
        userEvent.type(input, 'invalid-url');
        
        waitFor(() => {
          expect(screen.getByText('Invalid URL format.')).toBeInTheDocument();
        });
      });
    });
    describe('when url is valid', () => {
      it('should not show error message', async () => {
        render(
          <ConfigEditor
            onOptionsChange={onOptionsChange}
            options={testOptions}
          />
        );

        const input = screen.getByRole('textbox', { name: /url/i });
        userEvent.type(input, 'http://test.server.com:8080');

        waitFor(() => {
          expect(screen.queryByText('Invalid URL format.')).not.toBeInTheDocument();
        });
      });
      describe('when url is valid', () => {
        it('should not show error message', async () => {
          render(
            <ConfigEditor
              onOptionsChange={onOptionsChange}
              options={testOptions}
            />
          );
  
          const input = screen.getByRole('textbox', { name: /url/i });
          userEvent.type(input, 'https://a');
  
          waitFor(() => {
            expect(screen.queryByText('Invalid URL format.')).not.toBeInTheDocument();
          });
        });
      });
    });
  });
});
