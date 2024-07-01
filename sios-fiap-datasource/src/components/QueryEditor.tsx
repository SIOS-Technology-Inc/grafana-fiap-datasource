import React ,{ useEffect } from 'react';
import { useForm, useFieldArray, Controller } from 'react-hook-form';

import { InlineFieldRow, InlineField, Input, Button, Checkbox, RadioButtonList } from '@grafana/ui';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from '../datasource';
import { MyDataSourceOptions, MyQuery } from '../types';

import { css } from '@emotion/css';

type Props = QueryEditorProps<DataSource, MyQuery, MyDataSourceOptions>;

const dateTimeFormat = /^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$/;
const dateFormat = /^\d{4}-\d{2}-\d{2}$/;

const isValidDateTime = (value: string) => {
  if (value === '') {
    return true;
  }

  if (!dateTimeFormat.test(value) && !dateFormat.test(value)) {
    return "Invalid date format. Please use 'YYYY-MM-DD HH:MM:SS' or 'YYYY-MM-DD' format.";
  }

  const dateTime = new Date(value);
  if (isNaN(dateTime.getTime())) {
    return "Invalid date. Please check the values";
  }
  
  return true;
};

const pointValidationRule = {
  required: 'This field is required',
}

const handleBlur = (e: React.FocusEvent<HTMLInputElement>, field: any, query: MyQuery ,linkDashBords: string, timeKey: string,onChange: (value: any) => void) => {
  const inputValue = e.currentTarget.value;

  if(dateFormat.test(inputValue)) {
    const newValue = `${inputValue} 00:00:00`;
    field.onChange(newValue);
  }
}

export function QueryEditor({ query, onChange }: Props) {
  const fieldLimit = 100;

  const { control, watch, trigger } = useForm<MyQuery>({
    mode: 'onChange',
    reValidateMode: 'onChange',
    defaultValues: {
      point_ids: query.point_ids,
      data_range: query.data_range,
      start_time: {time: query.start_time.time, link_dashboard: query.start_time.link_dashboard},
      end_time: {time: query.end_time.time, link_dashboard: query.end_time.link_dashboard},
    }
  });

  const pointIds = watch('point_ids');
  const startLinkDashboards = watch('start_time.link_dashboard');
  const endLinkDashboards = watch('end_time.link_dashboard');
  const startLinkDashboardsValue = watch('start_time.time');
  const endLinkDashboardsValue = watch('end_time.time');
  
  const { fields, append,remove } = useFieldArray({
    name: 'point_ids',
    control,
  });
  
  useEffect(() => {
    trigger();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  },[pointIds])

  return (
    <>
      {fields.map((field, index) => {
        return(
          <div key={field.id}>
            <InlineFieldRow>
              <Controller
                name={`point_ids.${index}.point_id`}
                control={control}
                rules={pointValidationRule}
                render={({ field , fieldState:{ error } }) => (
                  <InlineField label="point" labelWidth={16} invalid={Boolean(error)} error={error?.message}>
                    <Input
                      id={`point-${index}`}
                      width={52}
                      placeholder="https://~~~"
                      value={field.value}
                      onChange={(e) => {
                        field.onChange(e);
                        onChange({ ...query, point_ids: pointIds})
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
                    onChange({ ...query, point_ids: query.point_ids.filter((_, i) => i !== index)})
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
                name={`data_range_${query.refId}`}
                options={[
                  { label: 'Period', value: 'period' },
                  { label: 'Latest', value: 'latest' },
                  { label: 'Oldest', value: 'oldest' }
                ]}
                value={field.value}
                onChange={(value) => {
                  field.onChange(value);
                  onChange({ ...query, data_range: value });
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
        <Controller
          name="start_time.time"
          control={control}
          rules={{ validate: isValidDateTime }}
          render={({ field, fieldState: { error } }) => (
            <InlineField label="Start" labelWidth={16} invalid={Boolean(error)} error={error && error.message}>
            <div style={{ pointerEvents: startLinkDashboards ? 'none' : 'auto', opacity: startLinkDashboards ? 0.4 : 1 }}>
              <Input
                id={`start_time`}
                placeholder="YYYY-MM-DD HH:MM:SS"
                value={field.value}
                onChange={(e) => {
                  field.onChange(e.currentTarget.value);
                  onChange({ ...query, start_time: { time: e.currentTarget.value, link_dashboard: startLinkDashboards } });
                }}
                onBlur={(e) => handleBlur(e, field, query, startLinkDashboardsValue , 'start_time', onChange)}
              />
            </div>
            </InlineField>
          )}
        />
        <Controller
          name="start_time.link_dashboard"
          control={control}
          render={({ field }) => (
            <Checkbox
              label='sync with grafana start time'
              onChange={(e) => {
                field.onChange(e.currentTarget.checked);
                onChange({ ...query, start_time: { time: startLinkDashboardsValue, link_dashboard: e.currentTarget.checked } });
              }}
                checked={field.value}
              />
              )}
            />
            </InlineFieldRow>
            <InlineFieldRow>
              <Controller
                name="end_time.time"
                control={control}
                rules={{ validate: isValidDateTime }}
                render={({ field ,fieldState:{ error }}) => (
                  <InlineField label="End" labelWidth={16} invalid={Boolean(error)} error={error && error.message}>
                  <div style={{ pointerEvents: endLinkDashboards ? 'none' : 'auto', opacity: endLinkDashboards ? 0.4 : 1 }}>
                  <Input
                    id={`end_time`}
                    placeholder="YYYY-MM-DD HH:MM:SS"
                    value={field.value}
                    onChange={(e) => {
                      field.onChange(e.currentTarget.value);
                      onChange({ ...query, end_time: { time: e.currentTarget.value, link_dashboard: endLinkDashboards } });
                    }}
                    onBlur={(e) => handleBlur(e, field, query, endLinkDashboardsValue , 'end_time', onChange)}
                  />
                  </div>
                  </InlineField>
                )}
              />
            <Controller
              name="end_time.link_dashboard"
              control={control}
              render={({ field }) => (
              <Checkbox
                label='sync with grafana end time'
                onChange={(e) => {
                  field.onChange(e.currentTarget.checked);
                  onChange({ ...query, end_time: { time: endLinkDashboardsValue, link_dashboard: e.currentTarget.checked } });
                }}
              checked={field.value}
            />
          )}
        />
      </InlineFieldRow>
    </>
  );
}
