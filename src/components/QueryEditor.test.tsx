import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { userEvent } from '@testing-library/user-event';

import { QueryEditor } from './QueryEditor';
import { MyQuery } from 'types';

describe('QueryEditor', () => {
  const testQuery: MyQuery = {
    "point_ids": [
      {
        "point_id": ""
      }
    ],
    "data_range": "period",
    "start_time": {
      "time": "",
      "link_dashboard": true
    },
    "end_time": {
      "time": "",
      "link_dashboard": true
    },
    "datasource": {
      "type": "sios-fiap-datasource",
      "uid": "a25326a4-956f-4006-8b21-1e4c3544772d"
    },
    "refId": "A"
  }
  const testDatasource = {
    "name": "Fiap",
    "id": 2,
    "type": "sios-fiap-datasource",
    "meta": {
      "id": "sios-fiap-datasource",
      "type": "datasource",
      "name": "Fiap",
      "info": {
        "author": {
          "name": "Sios",
          "url": ""
        },
        "description": "This is a grafana data source that uses ieee1888",
        "links": [],
        "logos": {
          "small": "public/plugins/sios-fiap-datasource/img/logo.svg",
          "large": "public/plugins/sios-fiap-datasource/img/logo.svg"
        },
        "build": {},
        "screenshots": [],
        "version": "1.0.0",
        "updated": "2024-06-20"
      },
      "dependencies": {
        "grafanaDependency": ">=10.0.3",
        "grafanaVersion": "*",
        "plugins": []
      },
      "includes": null,
      "category": "",
      "preload": false,
      "backend": true,
      "routes": null,
      "skipDataQuery": false,
      "autoEnabled": false,
      "annotations": false,
      "metrics": true,
      "alerting": false,
      "explore": false,
      "tables": false,
      "logs": false,
      "tracing": false,
      "streaming": false,
      "executable": "gpx_fiap",
      "signature": "unsigned",
      "module": "plugins/sios-fiap-datasource/module",
      "baseUrl": "public/plugins/sios-fiap-datasource"
    },
    "cachingConfig": {
      "enabled": false,
      "TTLMs": 0
    },
    "uid": "a25326a4-956f-4006-8b21-1e4c3544772d",
    "components": {}
  } as any;
  const testOnChange = jest.fn();
  const testonRunQuery = jest.fn();

  describe('point_ids validation test', () => {
    describe('when point_ids is empty in the initial state', () => {
      it('should show error message', async () => {
        render(<QueryEditor query={testQuery} onChange={testOnChange} onRunQuery={testonRunQuery} datasource={testDatasource}/>);

        await waitFor(() => {
          expect(screen.queryByText('This field is required')).toBeInTheDocument();
        });
      });
    });
    describe('when point_ids is empty after input', () => {
      it('should show error message', async () => {
        render(<QueryEditor query={testQuery} onChange={testOnChange} onRunQuery={testonRunQuery} datasource={testDatasource}/>);

        const input = screen.getByRole('textbox', { name: /point/i })
        await userEvent.type(input, 'a');
        await userEvent.clear(input);

        await waitFor(() => {
          expect(screen.queryByText('This field is required')).toBeInTheDocument();
        });
      });
    });
    describe('when one point_ids is filled', () => {
      it('should not show error message', async () => {
        render(<QueryEditor query={testQuery} onChange={testOnChange} onRunQuery={testonRunQuery} datasource={testDatasource}/>);

        const input = screen.getByRole('textbox', { name: /point/i })
        await userEvent.type(input, 'a');

        await waitFor(() => {
          expect(screen.queryByText('This field is required')).not.toBeInTheDocument();
        });
      });
    });
    describe('when two point_ids are empty after append one point', () => {
      it('should show two error messages', async () => {
        render(<QueryEditor query={testQuery} onChange={testOnChange} onRunQuery={testonRunQuery} datasource={testDatasource}/>);

        await userEvent.click(screen.getByRole('button', { name: 'plus' }));

        await waitFor(() => {
          const matchingErrors = screen.queryAllByText('This field is required');
          expect(matchingErrors).toHaveLength(2);
        });
      });
    });
    describe('when first point_id is filled and after append one point', () => {
      it('should show one error message', async () => {
        render(<QueryEditor query={testQuery} onChange={testOnChange} onRunQuery={testonRunQuery} datasource={testDatasource}/>);

        await userEvent.click(screen.getByRole('button', { name: 'plus' }));
        await userEvent.type(screen.getByTestId('point-0'), 'a');

        await waitFor(() => {
          const matchingErrors = screen.queryAllByText('This field is required');
          expect(matchingErrors).toHaveLength(1);
        });
      });
    });
    describe('when second point_id is filled and after append one point', () => {
      it('should show one error message', async () => {
        render(<QueryEditor query={testQuery} onChange={testOnChange} onRunQuery={testonRunQuery} datasource={testDatasource}/>);

        await userEvent.click(screen.getByRole('button', { name: 'plus' }));
        await userEvent.type(screen.getByTestId('point-1'), 'a');

        await waitFor(() => {
          const matchingErrors = screen.queryAllByText('This field is required');
          expect(matchingErrors).toHaveLength(1);
        });
      });
    });
    describe('when two point_ids are filled after append one point', () => {
      it('should not show error message', async () => {
        render(<QueryEditor query={testQuery} onChange={testOnChange} onRunQuery={testonRunQuery} datasource={testDatasource}/>);

        await userEvent.click(screen.getByRole('button', { name: 'plus' }));
        await userEvent.type(screen.getByTestId('point-0'), 'a');
        await userEvent.type(screen.getByTestId('point-1'), 'b');

        await waitFor(() => {
          expect(screen.queryByText('This field is required')).not.toBeInTheDocument();
        });
      });
    });
  });
  describe('time validation test', () => {
    const invalidFormatTimeInputs = [
      ' 2000-01-01 00:00:00',
      '2000-01-01 00:00:00 ',
      '20000-01-01 00:00:00',
      '2000-001-01 00:00:00',
      '2000-01-001 00:00:00',
      '2000-01-01 000:00:00',
      '2000-01-01 00:000:00',
      '2000-01-01 00:00:000',
      '200a-01-01 00:00:00',
      '2000-a1-01 00:00:00',
      '2000-01-a1 00:00:00',
      '2000-01-01 a0:00:00',
      '2000-01-01 00:a0:00',
      '2000-01-01 00:00:a0',
      '2000/01/01 00:00:00',
      '2000-01-01 00-00-00'
    ];
    const invalidTimeInputs = ['2000-13-01 00:00:00', '2000-12-32'];
    const validTimeInputs = ['2000-01-01 00:00:00', '2000-01-01'];

    describe('when start time and end time is empty in the initial state', () => {
      it('should not show date error message', async () => {
        render(<QueryEditor query={testQuery} onChange={testOnChange} onRunQuery={testonRunQuery} datasource={testDatasource}/>);

        await waitFor(() => {
          expect(screen.queryByText("Invalid date")).not.toBeInTheDocument();
        });
      });
    });
    describe('start time validation test', () => {
      describe('when start time format is invalid', () => {
        it.each(invalidFormatTimeInputs)('should show error message (input: %s)', async (invalidFormatTimeInput: string) => {
          render(<QueryEditor query={testQuery} onChange={testOnChange} onRunQuery={testonRunQuery} datasource={testDatasource}/>);
  
          await userEvent.click(screen.getByRole('checkbox', { name: /start/i }));
          await userEvent.type(screen.getByTestId('start-time-input'), invalidFormatTimeInput);
  
          await waitFor(() => {
            expect(screen.queryByText("Invalid date format. Please use 'YYYY-MM-DD HH:MM:SS' or 'YYYY-MM-DD' format.")).toBeInTheDocument();
          });
        });
      });
      describe('when start time format is valid but start time is invalid', () => {
        it.each(invalidTimeInputs)('should show error message (input:%s)', async (invalidTimeInput: string) => {
          render(<QueryEditor query={testQuery} onChange={testOnChange} onRunQuery={testonRunQuery} datasource={testDatasource}/>);
  
          await userEvent.click(screen.getByRole('checkbox', { name: /start/i }));
          await userEvent.type(screen.getByTestId('start-time-input'), invalidTimeInput);
  
          await waitFor(() => {
            expect(screen.queryByText("Invalid date. Please check the values")).toBeInTheDocument();
          });
        });
      });
      describe('when start time format and date is valid YYYY-MM-DD HH:MM:SS', () => {
        it.each(validTimeInputs)('should not show time error message (input: %s)', async (validTimeInput: string) => {
          render(<QueryEditor query={testQuery} onChange={testOnChange} onRunQuery={testonRunQuery} datasource={testDatasource}/>);
  
          await userEvent.click(screen.getByRole('checkbox', { name: /start/i }));
          await userEvent.type(screen.getByTestId('start-time-input'), validTimeInput);
  
          await waitFor(() => {
            expect(screen.queryByText("Invalid date")).not.toBeInTheDocument();
          });
        });
      });
    }); 
    describe('end time validation test', () => {
      describe('when end time format is invalid', () => {
        it.each(invalidFormatTimeInputs)('should show error message (input: %s)', async (invalidFormatTimeInput: string) => {
          render(<QueryEditor query={testQuery} onChange={testOnChange} onRunQuery={testonRunQuery} datasource={testDatasource}/>);
  
          await userEvent.click(screen.getByRole('checkbox', { name: /end/i }));
          await userEvent.type(screen.getByTestId('end-time-input'), invalidFormatTimeInput);
  
          await waitFor(() => {
            expect(screen.queryByText("Invalid date format. Please use 'YYYY-MM-DD HH:MM:SS' or 'YYYY-MM-DD' format.")).toBeInTheDocument();
          });
        });
      });
      describe('when end time format is valid but end time is invalid', () => {
        it.each(invalidTimeInputs)('should show error message (input:%s)', async (invalidTimeInput: string) => {
          render(<QueryEditor query={testQuery} onChange={testOnChange} onRunQuery={testonRunQuery} datasource={testDatasource}/>);
  
          await userEvent.click(screen.getByRole('checkbox', { name: /end/i }));
          await userEvent.type(screen.getByTestId('end-time-input'), invalidTimeInput);
  
          await waitFor(() => {
            expect(screen.queryByText("Invalid date. Please check the values")).toBeInTheDocument();
          });
        });
      });
      describe('when end time format and date is valid YYYY-MM-DD HH:MM:SS', () => {
        it.each(validTimeInputs)('should not show time error message (input: %s)', async (validTimeInput: string) => {
          render(<QueryEditor query={testQuery} onChange={testOnChange} onRunQuery={testonRunQuery} datasource={testDatasource}/>);
  
          await userEvent.click(screen.getByRole('checkbox', { name: /end/i }));
          await userEvent.type(screen.getByTestId('end-time-input'), validTimeInput);
  
          await waitFor(() => {
            expect(screen.queryByText("Invalid date")).not.toBeInTheDocument();
          });
        });
      });
    });
  });
  describe('custom logic test', () => {
    describe('when onBulr start time input after input YYYY-MM-DD format date', () => {
      it('should change start time format to YYYY-MM-DD 00:00:00', async () => {
        render(<QueryEditor query={testQuery} onChange={testOnChange} onRunQuery={testonRunQuery} datasource={testDatasource}/>);

        await userEvent.click(screen.getByRole('checkbox', { name: /start/i }));
        await userEvent.type(screen.getByTestId('start-time-input'), '2022-01-01');

        await userEvent.click(document.body);

        await waitFor(() => {
          expect(screen.queryByTestId('start-time-input')).toHaveValue('2022-01-01 00:00:00');
        });
      });
    });
    describe('when onBulr start time input after input YYYY-MM-DD HH:MM:SS format date', () => {
      it('should not change start time format', async () => {
        render(<QueryEditor query={testQuery} onChange={testOnChange} onRunQuery={testonRunQuery} datasource={testDatasource}/>);

        await userEvent.click(screen.getByRole('checkbox', { name: /start/i }));
        await userEvent.type(screen.getByTestId('start-time-input'), '2022-01-01 00:00:00');

        await userEvent.click(document.body);

        await waitFor(() => {
          expect(screen.queryByTestId('start-time-input')).toHaveValue('2022-01-01 00:00:00');
        });
      });
    });
    describe('when onBulr end time input after input YYYY-MM-DD format date', () => {
      it('should change end time format to YYYY-MM-DD 00:00:00', async () => {
        render(<QueryEditor query={testQuery} onChange={testOnChange} onRunQuery={testonRunQuery} datasource={testDatasource}/>);

        await userEvent.click(screen.getByRole('checkbox', { name: /end/i }));
        await userEvent.type(screen.getByTestId('end-time-input'), '2022-13-32');

        await userEvent.click(document.body);

        await waitFor(() => {
          expect(screen.queryByTestId('end-time-input')).toHaveValue('2022-13-32 00:00:00');
        });
      });
    });
    describe('when onBulr end time input after input YYYY-MM-DD HH:MM:SS format date', () => {
      it('should not change end time format', async () => {
        render(<QueryEditor query={testQuery} onChange={testOnChange} onRunQuery={testonRunQuery} datasource={testDatasource}/>);

        await userEvent.click(screen.getByRole('checkbox', { name: /end/i }));
        await userEvent.type(screen.getByTestId('end-time-input'), '2022-01-01 00:00:00');

        await userEvent.click(document.body);

        await waitFor(() => {
          expect(screen.queryByTestId('end-time-input')).toHaveValue('2022-01-01 00:00:00');
        });
      })
    });
  });
});