// Settings Manager - 设置管理器
// 提供设置页面的核心功能和状态管理

class SettingsManager {
    constructor() {
        this.initialized = false;
        this.currentConfig = null;
        this.init();
    }

    init() {
        if (this.initialized) return;
        
        console.log('Settings Manager initialized');
        this.initialized = true;
    }

    // 获取当前配置
    getCurrentConfig() {
        return this.currentConfig;
    }

    // 设置当前配置
    setCurrentConfig(config) {
        this.currentConfig = config;
    }

    // 验证配置
    validateConfig(config) {
        const errors = [];
        const warnings = [];

        try {
            // 基本验证
            if (!config) {
                errors.push('配置对象为空');
                return { errors, warnings };
            }

            // 验证应用配置
            if (!config.app?.title) {
                warnings.push('应用标题未设置');
            }

            // 验证课程配置
            if (!config.lessons?.schedule?.length) {
                warnings.push('课程表为空');
            }

            // 验证时间配置
            if (!config.lessons?.times?.schedule?.length) {
                warnings.push('时间表为空');
            }

            return { errors, warnings };
        } catch (error) {
            errors.push('配置验证失败: ' + error.message);
            return { errors, warnings };
        }
    }

    // 导出配置
    exportConfig() {
        try {
            if (typeof window.exportConfig === 'function') {
                return window.exportConfig();
            }
            throw new Error('导出功能不可用');
        } catch (error) {
            console.error('导出配置失败:', error);
            throw error;
        }
    }

    // 导入配置
    importConfig(file) {
        try {
            if (typeof window.importConfig === 'function') {
                return window.importConfig(file);
            }
            throw new Error('导入功能不可用');
        } catch (error) {
            console.error('导入配置失败:', error);
            throw error;
        }
    }

    // 重置配置
    resetConfig() {
        try {
            if (typeof window.resetConfig === 'function') {
                return window.resetConfig();
            }
            throw new Error('重置功能不可用');
        } catch (error) {
            console.error('重置配置失败:', error);
            throw error;
        }
    }

    // 保存配置
    saveConfig() {
        try {
            if (typeof window.handleSave === 'function') {
                return window.handleSave();
            }
            throw new Error('保存功能不可用');
        } catch (error) {
            console.error('保存配置失败:', error);
            throw error;
        }
    }
}

// 导出给全局使用
window.SettingsManager = SettingsManager;