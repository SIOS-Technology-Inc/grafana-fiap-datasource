import React, { FormEvent, useEffect } from 'react';
import { useForm, Controller } from 'react-hook-form';

import { InlineField, Input } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps, DataSourceSettings } from '@grafana/data';
import { MyDataSourceOptions } from '../types';

interface Props extends DataSourcePluginOptionsEditorProps<MyDataSourceOptions> {}

export const URL = '';

export function ConfigEditor(props: Props) {
  const { onOptionsChange, options } = props;

  const onURLChange = (event: FormEvent<HTMLInputElement>) => {
    const jsonData = {
      ...options.jsonData,
      url: event.currentTarget.value,
    };
    onOptionsChange({ ...options, jsonData });
  };

  const onJsonDataReset = (options: DataSourceSettings<MyDataSourceOptions, {}>) => {
    const tempJsonData = { ...options.jsonData };
    const jsonData = {
      ...tempJsonData,
      url: tempJsonData.url || URL,
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
      url: jsonData.url || URL,
    },
  });

  const urlValidationRule = {
    required: 'This field is required',
    pattern: {
      // (http|https):// に続いて半角英数字又は"_/:%#$&?()~.=+-"が1文字以上あることを確認する正規表現
      value: /^https?:\/\/[\w/:%#\$&\?\(\)~\.=\+\-]+$/,
      message: 'Invalid URL format.',
    },
  }

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
        rules={urlValidationRule}
        render={({ field, fieldState:{ error } }) => (
          <InlineField label="URL" labelWidth={12} invalid={Boolean(error)} error={error?.message}>
            <Input
              id='url'
              onChange={(e) => {
                field.onChange(e);
                onURLChange(e);
              }}
              value={field.value}
              placeholder="http://test.server.com:8080"
              width={40}
            />
          </InlineField>
        )}
      />
    </div>
  );
}
