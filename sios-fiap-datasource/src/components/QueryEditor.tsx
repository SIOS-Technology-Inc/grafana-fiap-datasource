import React , { useState } from 'react';
import { useForm, useFieldArray, Controller } from 'react-hook-form';

import { InlineFieldRow, InlineField, Input, Button, DateTimePicker, Checkbox, RadioButtonList } from '@grafana/ui';
import { DateTime, QueryEditorProps, dateTime } from '@grafana/data';
import { DataSource } from '../datasource';
import { MyDataSourceOptions, MyQueryForm, MyQuery } from '../types';

import { css } from '@emotion/css';
import { now } from 'lodash';

type Props = QueryEditorProps<DataSource, MyQuery, MyDataSourceOptions>;

export function QueryEditor({ query, onChange, onRunQuery }: Props) {
  const fieldLimit = 100;

  const { control, watch } = useForm<MyQueryForm>({
    mode: 'onChange',
    reValidateMode: 'onChange',
    defaultValues: {
      point_ids: [{ point_id : ''}],
      data_range: 'period',
      start_time: {time: new Date().toISOString(), link_dashboards: false},
      end_time: {time: new Date().toISOString(), link_dashboards: false},
    }
  });

  const startLinkDashboards = watch('start_time.link_dashboards');
  const endLinkDashboards = watch('end_time.link_dashboards');
  
  const { fields, append,remove } = useFieldArray({
    name: 'point_ids',
    control,
  });

  const [startDate, setStartDate] = useState<DateTime>(dateTime(now()));
  const [endDate, setEndDate] = useState<DateTime>(dateTime(now()));

  return (
    <>
      {fields.map((field, index) => {
        return(
          <div key={field.id}>
            <InlineFieldRow>
              <Controller
                name={`point_ids.${index}.point_id`}
                control={control}
                render={({ field }) => (
                  <InlineField label="point" labelWidth={16}>
                    <Input
                      id={`point-${index}`}
                      width={52}
                      placeholder="https://~~~"
                      value={field.value}
                      onChange={(e) => {
                        field.onChange(e);
                      }}
                    />
                  </InlineField>
                )}
              />
              {fields.length > 1 && (
                <Button
                  variant='secondary'
                  onClick={() => {
                    remove(index);
                  }}
                  title='minus'
                  icon='minus'
                />
              )}
              {index === fields.length -1 && fields.length < fieldLimit && (
                <Button
                  variant='secondary'
                  onClick={() => {
                    append({ point_id: '' });
                  }}
                  title='plus'
                  icon='plus'
                />
              )}
            </InlineFieldRow>
          </div>
        );
      })}
      <InlineFieldRow>
        <InlineField label="Data Range" labelWidth={16}>
          <Controller
            name="data_range"
            control={control}
            render={({ field }) => (
              <RadioButtonList
                name="date_range"
                options={[
                  { label: 'Period', value: 'period' },
                  { label: 'Latest', value: 'latest' },
                  { label: 'Oldest', value: 'oldest' }
                ]}
                value={field.value}
                onChange={(value) => {
                  field.onChange(value);
                } 
                }
                className={css`
                  grid-template-columns: 1fr 1fr 1fr;
                `}
              />
            )}
          />
        </InlineField>
      </InlineFieldRow>
      <InlineFieldRow>
        <InlineField label="Start" labelWidth={16}>
        <Controller
          name="end_time.time"
          control={control}
          render={({ field }) => (
            <div style={{ pointerEvents: startLinkDashboards ? 'none' : 'auto', opacity: startLinkDashboards ? 0.4 : 1 }}>
              <DateTimePicker
                date={startDate}
                onChange={(time) => {
                  setStartDate(time);
                  field.onChange(time.toISOString());
                }}
              />
            </div>
          )}
        />
        </InlineField>
        <Controller
          name="start_time.link_dashboards"
          control={control}
          render={({ field }) => (
            <Checkbox
              label='Grafanaの時間指定と連動'
              onChange={(e) => {
                field.onChange(e.currentTarget.checked);
                }}
                checked={field.value}
              />
              )}
            />
            </InlineFieldRow>
            <InlineFieldRow>
            <InlineField label="End" labelWidth={16}>
              <Controller
              name="end_time.time"
              control={control}
              render={({ field }) => (
                <div style={{ pointerEvents: endLinkDashboards ? 'none' : 'auto', opacity: endLinkDashboards ? 0.4 : 1 }}>
                <DateTimePicker
                  date={endDate}
                  onChange={(time) => {
                    setEndDate(time);
                    field.onChange(time.toISOString());
                  }}
                />
                </div>
              )}
            />
            </InlineField>
            <Controller
              name="end_time.link_dashboards"
              control={control}
              render={({ field }) => (
              <Checkbox
                label='Grafanaの時間指定と連動'
                onChange={(e) => {
                field.onChange(e.currentTarget.checked);
              }}
              checked={field.value}
            />
          )}
        />
      </InlineFieldRow>
    </>
  );
}
