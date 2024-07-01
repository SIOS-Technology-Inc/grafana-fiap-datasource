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

        waitFor(() => {
          expect(screen.getByText('This field is required')).toBeInTheDocument();
        });
      });
    });
    describe('when point_ids is empty after input', () => {
      it('should show error message', async () => {
        render(<QueryEditor query={testQuery} onChange={testOnChange} onRunQuery={testonRunQuery} datasource={testDatasource}/>);

        const input = screen.getByRole('textbox', { name: /point/i })
        userEvent.type(input, 'test point');
        userEvent.clear(input);

        waitFor(() => {
          expect(screen.getByText('This field is required')).toBeInTheDocument();
        });
      });
    });
    describe('when two point_ids are empty after append one point', () => {
      it('should show two error messages', async () => {
        render(<QueryEditor query={testQuery} onChange={testOnChange} onRunQuery={testonRunQuery} datasource={testDatasource}/>);

        userEvent.click(screen.getByRole('button', { name: 'plus' }));

        waitFor(() => {
          const matchingErrors = screen.getAllByText('This field is required');
          expect(matchingErrors).toHaveLength(2);
        });
      });
    });
  });
  describe('time validation test', () => {
    describe('when start time and end time is empty in the initial state', () => {
      it('should not show date error message', async () => {
        render(<QueryEditor query={testQuery} onChange={testOnChange} onRunQuery={testonRunQuery} datasource={testDatasource}/>);

        waitFor(() => {
          expect(screen.getByText("Invalid date")).not.toBeInTheDocument();
        });
      });
    });
    describe('when start time format is invalid (ex.2000-01-01 00-00-00)', () => {
      it('should show error message', async () => {
        render(<QueryEditor query={testQuery} onChange={testOnChange} onRunQuery={testonRunQuery} datasource={testDatasource}/>);

        userEvent.click(screen.getByTestId('start-time-input'));
        userEvent.type(screen.getByTestId('start-time-input'), '2000-01-01 00-00-00');

        waitFor(() => {
          expect(screen.getByText("Invalid date format. Please use 'YYYY-MM-DD HH:MM:SS' or 'YYYY-MM-DD' format.")).toBeInTheDocument();
        });
      });
    });
    describe('when start time format is invalid (ex.2000/01/01)', () => {
      it('should show error message', async () => {
        render(<QueryEditor query={testQuery} onChange={testOnChange} onRunQuery={testonRunQuery} datasource={testDatasource}/>);

        userEvent.click(screen.getByTestId('start-time-input'));
        userEvent.type(screen.getByTestId('start-time-input'), '2000/01/01');

        waitFor(() => {
          expect(screen.getByText("Invalid date format. Please use 'YYYY-MM-DD HH:MM:SS' or 'YYYY-MM-DD' format.")).toBeInTheDocument();
        });
      });
    });
    describe('when start time format is valid(YYYY-MM-DD HH:MM:SS) but start time is invalid', () => {
      it('should show error message', async () => {
        render(<QueryEditor query={testQuery} onChange={testOnChange} onRunQuery={testonRunQuery} datasource={testDatasource}/>);

        userEvent.click(screen.getByRole('checkbox', { name: /start/i }));
        userEvent.type(screen.getByTestId('start-time-input'), '2022-13-01 00:00:00');

        waitFor(() => {
          expect(screen.getByText("Invalid date. Please check the values")).toBeInTheDocument();
        });
      });
    });
    describe('when start time format is valid(YYYY-MM-DD) but start time is invalid', () => {
      it('should show error message', async () => {
        render(<QueryEditor query={testQuery} onChange={testOnChange} onRunQuery={testonRunQuery} datasource={testDatasource}/>);

        userEvent.click(screen.getByRole('checkbox', { name: /start/i }));
        userEvent.type(screen.getByTestId('start-time-input'), '2022-12-32');
      });
    });
    describe('when start time format and date is valid YYYY-MM-DD HH:MM:SS', () => {
      it('should not show time error message', async () => {
        render(<QueryEditor query={testQuery} onChange={testOnChange} onRunQuery={testonRunQuery} datasource={testDatasource}/>);

        userEvent.click(screen.getByRole('checkbox', { name: /start/i }));
        userEvent.type(screen.getByTestId('start-time-input'), '2022-01-01 00:00:00');

        waitFor(() => {
          expect(screen.getByText("Invalid date")).not.toBeInTheDocument();
        });
      });
    });
    describe('when start time format and date is valid YYYY-MM-DD', () => {
      it('should not show time error message', async () => {
        render(<QueryEditor query={testQuery} onChange={testOnChange} onRunQuery={testonRunQuery} datasource={testDatasource}/>);

        userEvent.click(screen.getByRole('checkbox', { name: /start/i }));
        userEvent.type(screen.getByTestId('start-time-input'), '2022-01-01');

        waitFor(() => {
          expect(screen.getByText("Invalid date")).not.toBeInTheDocument();
        });
      });
    });
    describe('when end time and end time is empty in the initial state', () => {
      it('should not show date error message', async () => {
        render(<QueryEditor query={testQuery} onChange={testOnChange} onRunQuery={testonRunQuery} datasource={testDatasource}/>);

        waitFor(() => {
          expect(screen.getByText("Invalid date")).not.toBeInTheDocument();
        });
      });
    });
    describe('when end time format is invalid (ex.2000-01-01 00-00-000)', () => {
      it('should show error message', async () => {
        render(<QueryEditor query={testQuery} onChange={testOnChange} onRunQuery={testonRunQuery} datasource={testDatasource}/>);

        userEvent.click(screen.getByTestId('end-time-input'));
        userEvent.type(screen.getByTestId('end-time-input'), '2000-01-01 00-00-000');

        waitFor(() => {
          expect(screen.getByText("Invalid date format. Please use 'YYYY-MM-DD HH:MM:SS' or 'YYYY-MM-DD' format.")).toBeInTheDocument();
        });
      });
    });
    describe('when end time format is invalid (ex.2000-01-01 )', () => {
      it('should show error message', async () => {
        render(<QueryEditor query={testQuery} onChange={testOnChange} onRunQuery={testonRunQuery} datasource={testDatasource}/>);

        userEvent.click(screen.getByTestId('end-time-input'));
        userEvent.type(screen.getByTestId('end-time-input'), '2000-01-01 ');

        waitFor(() => {
          expect(screen.getByText("Invalid date format. Please use 'YYYY-MM-DD HH:MM:SS' or 'YYYY-MM-DD' format.")).toBeInTheDocument();
        });
      });
    });
    describe('when end time format is valid(YYYY-MM-DD HH:MM:SS) but end time is invalid', () => {
      it('should show error message', async () => {
        render(<QueryEditor query={testQuery} onChange={testOnChange} onRunQuery={testonRunQuery} datasource={testDatasource}/>);

        userEvent.click(screen.getByRole('checkbox', { name: /end/i }));
        userEvent.type(screen.getByTestId('end-time-input'), '2022-13-01 00:00:00');

        waitFor(() => {
          expect(screen.getByText("Invalid date. Please check the values")).toBeInTheDocument();
        });
      });
    });
    describe('when end time format is valid(YYYY-MM-DD) but end time is invalid', () => {
      it('should show error message', async () => {
        render(<QueryEditor query={testQuery} onChange={testOnChange} onRunQuery={testonRunQuery} datasource={testDatasource}/>);

        userEvent.click(screen.getByRole('checkbox', { name: /end/i }));
        userEvent.type(screen.getByTestId('end-time-input'), '2022-12-32');
      });
    });
    describe('when end time format and date is valid YYYY-MM-DD HH:MM:SS', () => {
      it('should not show time error message', async () => {
        render(<QueryEditor query={testQuery} onChange={testOnChange} onRunQuery={testonRunQuery} datasource={testDatasource}/>);

        userEvent.click(screen.getByRole('checkbox', { name: /end/i }));
        userEvent.type(screen.getByTestId('end-time-input'), '2022-01-01 00:00:00');

        waitFor(() => {
          expect(screen.getByText("Invalid date")).not.toBeInTheDocument();
        });
      });
    });
    describe('when end time format and date is valid YYYY-MM-DD', () => {
      it('should not show time error message', async () => {
        render(<QueryEditor query={testQuery} onChange={testOnChange} onRunQuery={testonRunQuery} datasource={testDatasource}/>);

        userEvent.click(screen.getByRole('checkbox', { name: /end/i }));
        userEvent.type(screen.getByTestId('end-time-input'), '2022-01-01');

        waitFor(() => {
          expect(screen.getByText("Invalid date")).not.toBeInTheDocument();
        });
      });
    });
  });
  describe('custom logic test', () => {
    describe('when onBulr start time input after input YYYY-MM-DD format date', () => {
      it('should change start time format to YYYY-MM-DD 00:00:00', async () => {
        render(<QueryEditor query={testQuery} onChange={testOnChange} onRunQuery={testonRunQuery} datasource={testDatasource}/>);

        userEvent.click(screen.getByRole('checkbox', { name: /start/i }));
        userEvent.type(screen.getByTestId('start-time-input'), '2022-01-01');

        userEvent.click(document.body);

        waitFor(() => {
          expect(screen.getByTestId('start-time-input')).toHaveValue('2022-01-01 00:00:00');
        });
      });
    });
    describe('when onBulr end time input after input YYYY-MM-DD format date', () => {
      it('should change end time format to YYYY-MM-DD 00:00:00', async () => {
        render(<QueryEditor query={testQuery} onChange={testOnChange} onRunQuery={testonRunQuery} datasource={testDatasource}/>);

        userEvent.click(screen.getByRole('checkbox', { name: /end/i }));
        userEvent.type(screen.getByTestId('start-time-input'), '2022-13-32');

        userEvent.click(document.body);

        waitFor(() => {
          expect(screen.getByTestId('start-time-input')).toHaveValue('2022-13-32 00:00:00');
        });
      });
    });
  });
});