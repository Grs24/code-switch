import { CurrentVersion, CheckForUpdates } from '../../bindings/codeswitch/versionservice'
import type { ReleaseInfo } from '../../bindings/codeswitch/models'

export const fetchCurrentVersion = async (): Promise<string> => {
  const version = await CurrentVersion()
  return version ?? ''
}

export const checkForUpdates = async (): Promise<ReleaseInfo | null> => {
  try {
    const release = await CheckForUpdates()
    return release
  } catch (error) {
    console.error('Failed to check for updates:', error)
    return null
  }
}
