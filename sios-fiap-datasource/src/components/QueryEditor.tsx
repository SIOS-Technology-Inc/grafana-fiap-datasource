import React , { useEffect } from 'react';
import { useForm, useFieldArray, Controller } from 'react-hook-form';

import { InlineFieldRow, InlineField, Input, Button, Checkbox, RadioButtonList } from '@grafana/ui';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from '../datasource';
import { MyDataSourceOptions, MyQueryForm, MyQuery } from '../types';

import { css } from '@emotion/css';

import { transformMyQueryFormToMyQuery } from '../utils/transformMyQueryFormToMyQuery';

type Props = QueryEditorProps<DataSource, MyQuery, MyDataSourceOptions>;

const isValidDateTime = (value: string) => {
  // 空文字を許容する
  if (value === '') {
    return true;
  }

  const dateTimeFormat = /^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$/;
  if (!dateTimeFormat.test(value)) {
    return "Invalid date format. Please use 'YYYY-MM-DD HH:MM:SS'";
  }
  const dateTime = new Date(value);
  if (isNaN(dateTime.getTime())) {
    return "Invalid date. Please check the values";
  }
  return true;
};


export function QueryEditor({ query, onChange, onRunQuery}: Props) {
  const fieldLimit = 100;

  const { control, watch } = useForm<MyQueryForm>({
    mode: 'onChange',
    reValidateMode: 'onChange',
    defaultValues: {
      point_ids: [{ point_id : ''}],
      data_range: 'period',
      start_time: {time: '', link_dashboard: false},
      end_time: {time: '', link_dashboard: false},
    }
  });

  const form = watch();

  useEffect(() => {
    const transformedQuery = transformMyQueryFormToMyQuery(form, query.refId);
    onChange(transformedQuery); // 変換後のqueryをonChangeで更新
  }, [form, query.refId, onChange]); // formが変更されるたびにこのeffectが実行されます

  const startLinkDashboards = watch('start_time.link_dashboard');
  const endLinkDashboards = watch('end_time.link_dashboard');
  
  const { fields, append,remove } = useFieldArray({
    name: 'point_ids',
    control,
  });

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
        <Controller
          name="start_time.time"
          control={control}
          rules={{ validate: isValidDateTime }}
          render={({ field, fieldState: { error } }) => (
            <InlineField label="Start" labelWidth={16} invalid={Boolean(error)} error={error && error.message}>
              <Input
                id={`start_time`}
                placeholder="YYYY-MM-DD HH:MM:SS"
                value={field.value}
                onChange={(e) => {
                  field.onChange(e.currentTarget.value);
                }}
                disabled={startLinkDashboards}
              />
            </InlineField>
          )}
        />
        <Controller
          name="start_time.link_dashboard"
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
              <Controller
                name="end_time.time"
                control={control}
                rules={{ validate: isValidDateTime }}
                render={({ field ,fieldState:{ error }}) => (
                  <InlineField label="End" labelWidth={16} invalid={Boolean(error)} error={error && error.message}>
                  <Input
                    id={`end_time`}
                    placeholder="YYYY-MM-DD HH:MM:SS"
                    value={field.value}
                    onChange={(e) => {
                      field.onChange(e.currentTarget.value);
                    }}
                    disabled={endLinkDashboards}
                  />
                  </InlineField>
                )}
              />
            <Controller
              name="end_time.link_dashboard"
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
