import { describe, expect, test } from 'bun:test'
import {
  CHANNEL_FORM_DEFAULT_VALUES,
  transformChannelToFormDefaults,
  transformFormDataToCreatePayload,
  transformFormDataToUpdatePayload,
  type ChannelFormValues,
} from './channel-form'
import type { Channel } from '../types'

function createChannel(overrides: Partial<Channel> = {}): Channel {
  return {
    id: 1,
    type: 1,
    key: '',
    openai_organization: '',
    test_model: '',
    status: 1,
    name: 'test',
    weight: 0,
    created_time: 0,
    test_time: 0,
    response_time: 0,
    base_url: '',
    other: '',
    balance: 0,
    balance_updated_time: 0,
    models: 'gpt-4o',
    group: 'default',
    used_quota: 0,
    model_mapping: '',
    status_code_mapping: '',
    priority: 0,
    auto_ban: 1,
    other_info: '',
    tag: '',
    setting: '',
    param_override: '',
    header_override: '',
    remark: '',
    max_input_tokens: 0,
    channel_info: {
      is_multi_key: false,
      multi_key_size: 0,
      multi_key_polling_index: 0,
      multi_key_mode: 'random',
    },
    settings: '{}',
    ...overrides,
  }
}

function createFormValues(
  overrides: Partial<ChannelFormValues> = {}
): ChannelFormValues {
  return {
    ...CHANNEL_FORM_DEFAULT_VALUES,
    name: 'test',
    models: 'gpt-4o',
    group: ['default'],
    ...overrides,
  }
}

function parseSetting(
  setting: string | null | undefined
): Record<string, unknown> {
  expect(setting).toBeTypeOf('string')
  return JSON.parse(setting || '{}')
}

describe('channel response_format form settings', () => {
  test('defaults unified response format to disabled for new channels', () => {
    expect(CHANNEL_FORM_DEFAULT_VALUES.response_format_enabled).toBe(false)

    const payload = transformFormDataToCreatePayload(createFormValues())
    const setting = parseSetting(payload.channel.setting)

    expect(setting.response_format).toEqual({
      enabled: false,
      mode: 'client_stream',
      rules: [],
    })
  })

  test('defaults legacy channels without response_format to disabled', () => {
    const formValues = transformChannelToFormDefaults(
      createChannel({
        setting: JSON.stringify({ force_format: true }),
      })
    )

    expect(formValues.response_format_enabled).toBe(false)
  })

  test('reads existing response_format enabled state when editing', () => {
    const formValues = transformChannelToFormDefaults(
      createChannel({
        setting: JSON.stringify({
          response_format: {
            enabled: true,
            mode: 'client_stream',
            rules: [{ future: 'kept' }],
          },
        }),
      })
    )

    expect(formValues.response_format_enabled).toBe(true)
  })

  test('writes enabled response_format and preserves unknown rule fields', () => {
    const payload = transformFormDataToUpdatePayload(
      createFormValues({
        response_format_enabled: true,
        setting: JSON.stringify({
          response_format: {
            enabled: false,
            mode: 'client_stream',
            rules: [
              {
                enabled: true,
                future_field: 'kept',
                nested: { value: 1 },
              },
            ],
          },
        }),
      }),
      1
    )
    const setting = parseSetting(payload.setting)

    expect(setting.response_format).toEqual({
      enabled: true,
      mode: 'client_stream',
      rules: [
        {
          enabled: true,
          future_field: 'kept',
          nested: { value: 1 },
        },
      ],
    })
  })

  test('writes a complete disabled response_format object', () => {
    const payload = transformFormDataToUpdatePayload(
      createFormValues({
        response_format_enabled: false,
        setting: JSON.stringify({
          response_format: {
            enabled: true,
            mode: 'client_stream',
            rules: [{ future_field: 'cleared' }],
          },
        }),
      }),
      1
    )
    const setting = parseSetting(payload.setting)

    expect(setting.response_format).toEqual({
      enabled: false,
      mode: 'client_stream',
      rules: [],
    })
  })
})
