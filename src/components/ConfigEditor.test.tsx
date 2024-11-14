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
      "url": "",
      "server_timezone": ""
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
  describe('server timezone validation test', () => {
    describe('when server timezone is empty in the initial state', () => {
      it('should not show error message', async () => {
        render(<ConfigEditor onOptionsChange={onOptionsChange} options={testOptions} />);

        await waitFor(() => {
          expect(screen.queryByText('Invalid timezone format. Please use the format ±HH:MM. For example, +09:00 or -05:30.')).not.toBeInTheDocument();
        });
      });
    });
    describe('when server timezone is invalid', () => {
      const InputTimezones = ['+0012', '-0012', '+13:00', '-13:00', '+00:60', '-00:60', '+00:0', '-00:0', '+0:00', '-0:00', '+0:0', '-0:0','+12','-12']
      it.each(InputTimezones)('should show error message (input: %s)', async (inputTimezone) => {
        render(<ConfigEditor onOptionsChange={onOptionsChange} options={testOptions} />);

        const input = screen.getByRole('textbox', { name: /server timezone/i });

        await userEvent.type(input, inputTimezone);
        
        await waitFor(() => {
          expect(screen.queryByText('Invalid timezone format. Please use the format ±HH:MM. For example, +09:00 or -05:30.')).toBeInTheDocument();
        });
      });
    });
    describe('when server timezone is valid', () => {
      const InputTimezones = ['+00:00', '-00:00', '+12:59', '-12:59', '+00:30', '-00:30', '+09:00', '-09:00', '+05:30', '-05:30']
      it.each(InputTimezones)('should not show error message (input: %s)', async (inputTimezone) => {
        render(<ConfigEditor onOptionsChange={onOptionsChange} options={testOptions} />);

        const input = screen.getByRole('textbox', { name: /server timezone/i });

        await userEvent.type(input, inputTimezone);
        
        await waitFor(() => {
          expect(screen.queryByText('Invalid timezone format. Please use the format ±HH:MM. For example, +09:00 or -05:30.')).not.toBeInTheDocument();
        });
      });
    });
  });
});
