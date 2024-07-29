import React, { useEffect } from 'react';
import { useForm, Controller } from 'react-hook-form';

import { InlineField, Input } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps, DataSourceSettings } from '@grafana/data';
import { MyDataSourceOptions } from '../types';

interface Props extends DataSourcePluginOptionsEditorProps<MyDataSourceOptions> {}

export function ConfigEditor(props: Props) {
  const { onOptionsChange, options} = props;

  const onJsonDataReset = (options: DataSourceSettings<MyDataSourceOptions, {}>) => {
    const tempJsonData = { ...options.jsonData };
    const jsonData = {
      ...tempJsonData,
      url: tempJsonData.url || '',
      server_timezone: tempJsonData.server_timezone || '',
    };
    onOptionsChange({ ...options, jsonData });
  };

  const { jsonData } = options;

  // React Hook Form
  const { control,trigger } = useForm<MyDataSourceOptions>({
    mode: 'onChange',
    reValidateMode: 'onChange',
    criteriaMode: 'all',
    defaultValues: {
      url: jsonData.url || '',
      server_timezone: jsonData.server_timezone || '',
    },
  });

  const validationRule = {
    url: {
      required: 'This field is required',
      pattern: {
        // (http|https):// に続いて1文字以上の文字列が続く形式
        value: /^https?:\/\/.+/,
        message: 'Invalid URL format.',
      },
    },
    server_timezone: {
      // server timezoneが入力されている場合のみチェック
      pattern: {
        // ±HH:MMの形式。±は+か-のどちらか、HHは00から12、MMは00から59の数字
        value: /^(\+|-)(0[0-9]|1[0-2]):[0-5][0-9]$/,
        message: 'Invalid timezone format. Please use the format ±HH:MM. For example, +09:00 or -05:30.',
      },
    },
  };

  useEffect(() => {
    onJsonDataReset(options);
    trigger();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  },[])

  return (
    <div className="gf-form-group">
      <Controller
        name="url"
        control={control}
        rules={validationRule.url}
        render={({ field, fieldState:{ error } }) => (
          <InlineField label="URL" labelWidth={18} invalid={Boolean(error)} error={error?.message}>
            <Input
              id='url'
              onChange={(e) => {
                field.onChange(e);
                onOptionsChange({ ...options, jsonData: { ...options.jsonData, url: e.currentTarget.value } });
              }}
              value={field.value}
              placeholder="http://test.server.com:8080"
              width={40}
            />
          </InlineField>
        )}
      />
      <Controller
        name="server_timezone"
        control={control}
        rules={validationRule.server_timezone}
        render={({ field, fieldState:{ error } }) => (
          <InlineField label="Server timezone" labelWidth={18} tooltip={"UTC is the default setting. If the field is empty, UTC will be used."} invalid={Boolean(error)} error={error?.message}>
            <Input
              id='server_timezone'
              onChange={(e) => {
                field.onChange(e);
                onOptionsChange({ ...options, jsonData: { ...options.jsonData, server_timezone: e.currentTarget.value } });
              }}
              value={field.value}
              placeholder="+09:00"
              width={10}
            />
          </InlineField>
        )}
      />
    </div>
  );
}
